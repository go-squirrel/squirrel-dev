# squctl 配置设计

## 概述

`squctl` 是命令行工具，模仿前端请求 APIServer，使用 Token 认证。

## 配置文件

`~/.squctl/config.yaml`：

```yaml
server: https://apiserver.example.com:10700
token: "eyJhbGciOiJIUzI1NiIs..."
```

## 命令

```bash
# 交互式登录（回车后依次输入）
squctl login
Server: https://apiserver.example.com:10700
Username: admin
Password: ****
Login Succeeded

# 部分免交互（指定 server，跳过 server 输入）
squctl login --apiserver https://apiserver.example.com:10700
Username: admin
Password: ****
Login Succeeded

# 完全免交互（全部指定，直接登录）
squctl login --apiserver https://... --user admin --password ****
Login Succeeded

# 登出
squctl logout

# 查看配置
squctl config view
```

## 实现

```go
type Config struct {
    Server string `yaml:"server"`
    Token  string `yaml:"token"`
}

func (c *Config) Save() error {
    data, _ := yaml.Marshal(c)
    return os.WriteFile(DefaultConfigPath(), data, 0600)
}
```

请求时携带 Token：

```go
req.Header.Set("Authorization", "Bearer "+config.Token)
```
