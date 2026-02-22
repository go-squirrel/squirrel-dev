# 证书生成指南【未验证】

本文档介绍如何使用 `cfssl` 工具手动生成 mTLS 双向认证所需的证书。

## 前置条件

安装 cfssl 工具：

```bash
# 下载 cfssl 和 cfssljson
curl -L -o /usr/local/bin/cfssl https://github.com/cloudflare/cfssl/releases/download/v1.6.4/cfssl_1.6.4_linux_amd64
curl -L -o /usr/local/bin/cfssljson https://github.com/cloudflare/cfssl/releases/download/v1.6.4/cfssljson_1.6.4_linux_amd64

# 添加执行权限
chmod +x /usr/local/bin/cfssl /usr/local/bin/cfssljson

# 验证安装
cfssl version
```

## 证书目录结构

```
certs/
├── ca.crt          # CA 证书
├── ca.key          # CA 私钥（妥善保管）
├── server.crt      # APIServer 服务端证书
├── server.key      # APIServer 服务端私钥
├── agent.crt       # Agent 客户端证书
└── agent.key       # Agent 客户端私钥
```

## 1. 创建证书目录

```bash
mkdir -p certs
cd certs
```

## 2. 生成 CA 证书

创建 CA 证书配置文件 `ca-config.json`：

```json
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "server": {
        "expiry": "87600h",
        "usages": ["signing", "key encipherment", "server auth"]
      },
      "client": {
        "expiry": "87600h",
        "usages": ["signing", "key encipherment", "client auth"]
      }
    }
  }
}
```

创建 CA 证书签名请求 `ca-csr.json`：

```json
{
  "CN": "squirrel-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "Beijing",
      "O": "Squirrel",
      "OU": "CA",
      "ST": "Beijing"
    }
  ]
}
```

生成 CA 证书：

```bash
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

生成文件：
- `ca.pem` → 重命名为 `ca.crt`
- `ca-key.pem` → 重命名为 `ca.key`

```bash
mv ca.pem ca.crt
mv ca-key.pem ca.key
```

## 3. 生成 APIServer 服务端证书

创建服务端证书签名请求 `server-csr.json`：

> **注意**：`hosts` 列表需要包含 APIServer 的所有访问地址（IP、域名、localhost）。

```json
{
  "CN": "squirrel-apiserver",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "hosts": [
    "127.0.0.1",
    "localhost",
    "squirrel-apiserver"
  ],
  "names": [
    {
      "C": "CN",
      "L": "Beijing",
      "O": "Squirrel",
      "OU": "APIServer",
      "ST": "Beijing"
    }
  ]
}
```

生成服务端证书：

```bash
cfssl gencert \
  -ca=ca.crt \
  -ca-key=ca.key \
  -config=ca-config.json \
  -profile=server \
  server-csr.json | cfssljson -bare server
```

生成文件：
- `server.pem` → 重命名为 `server.crt`
- `server-key.pem` → 重命名为 `server.key`

```bash
mv server.pem server.crt
mv server-key.pem server.key
```

## 4. 生成 Agent 客户端证书

创建客户端证书签名请求 `agent-csr.json`：

> **注意**：`CN` 字段用于标识 Agent 身份，可用于 ACL 控制。

```json
{
  "CN": "squirrel-agent",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "hosts": [],
  "names": [
    {
      "C": "CN",
      "L": "Beijing",
      "O": "Squirrel",
      "OU": "Agent",
      "ST": "Beijing"
    }
  ]
}
```

生成客户端证书：

```bash
cfssl gencert \
  -ca=ca.crt \
  -ca-key=ca.key \
  -config=ca-config.json \
  -profile=client \
  agent-csr.json | cfssljson -bare agent
```

生成文件：
- `agent.pem` → 重命名为 `agent.crt`
- `agent-key.pem` → 重命名为 `agent.key`

```bash
mv agent.pem agent.crt
mv agent-key.pem agent.key
```

## 5. 分发证书

### APIServer

APIServer 需要以下证书文件：

```bash
# 在 APIServer 节点上创建目录
mkdir -p /etc/squirrel/certs

# 复制证书
cp ca.crt server.crt server.key /etc/squirrel/certs/
```

### Agent

Agent 需要以下证书文件：

```bash
# 在 Agent 节点上创建目录
mkdir -p /etc/squirrel/certs

# 复制证书
cp ca.crt agent.crt agent.key /etc/squirrel/certs/
```

## 6. 配置 APIServer

修改 `config/apiserver.yaml`：

```yaml
mtls:
  enabled: true
  caFile: /etc/squirrel/certs/ca.crt
  certFile: /etc/squirrel/certs/server.crt
  keyFile: /etc/squirrel/certs/server.key
  allowedCNs:
    - squirrel-agent        # 只允许 squirrel-agent 连接
    # - squirrel-admin      # 可添加其他允许的 CN
```

## 7. 验证证书

验证证书内容：

```bash
# 查看 CA 证书
openssl x509 -in ca.crt -text -noout

# 查看服务端证书
openssl x509 -in server.crt -text -noout

# 查看客户端证书
openssl x509 -in agent.crt -text -noout

# 验证证书链
openssl verify -CAfile ca.crt server.crt
openssl verify -CAfile ca.crt agent.crt
```

测试 mTLS 连接：

```bash
# 使用客户端证书连接 APIServer
curl --cacert ca.crt \
     --cert agent.crt \
     --key agent.key \
     https://127.0.0.1:10700/api/v1/health
```

## 安全建议

1. **保护 CA 私钥**：`ca.key` 应妥善保管，不要分发到生产节点
2. **证书有效期**：默认 10 年，生产环境建议设置更短的有效期
3. **最小权限**：`allowedCNs` 只添加必要的客户端 CN
4. **定期轮换**：建议定期轮换证书，降低泄露风险
