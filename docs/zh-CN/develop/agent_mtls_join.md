# Agent mTLS 安全注册机制设计

## 背景

当前 Agent 与 APIServer 之间的通信存在以下安全问题：

1. **缺少身份验证**：Agent 与 APIServer 通信未使用双向认证
2. **证书分发困难**：需要手动在每台 Agent 机器上部署证书
3. **Token 注册机制缺失**：没有类似 `kubeadm join` 的安全加入机制

## 设计目标

1. **mTLS 双向认证**：Agent 与 APIServer 之间使用 mTLS 进行通信
2. **Token 自动注册**：类似 `kubeadm join`，使用临时 Token 完成 Agent 加入
3. **证书自动分发**：Agent 注册成功后，APIServer 自动签发并分发 Agent 专属证书

## 整体架构

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          Agent 安全注册流程                                   │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 生成 Join Token                                                          │
│  ┌─────────────┐                                                             │
│  │   Admin     │ ──► squctl token create ──► Token: abcdef.1234567890abcdef  │
│  │  (控制节点)  │                                             │              │
│  └─────────────┘                                             │              │
│                                                              ▼              │
│  2. Agent 加入                                              复制 Token      │
│  ┌─────────────┐                                            到 Agent       │
│  │   Agent     │ ◄─────────────────────────────────────────────────         │
│  │  (工作节点)  │                                                             │
│  └──────┬──────┘                                                             │
│         │                                                                    │
│         │ POST /api/v1/agent/join                                            │
│         │ { token: "abcdef.1234567890abcdef", hostname, ip, ... }            │
│         │ [使用 Bootstrap 证书 或 仅 HTTP]                                    │
│         ▼                                                                    │
│  3. APIServer 验证 Token                                                      │
│  ┌─────────────────┐                                                         │
│  │   APIServer     │ ──► 验证 Token 有效性（未过期、未使用）                    │
│  │                 │ ──► 生成 Agent 专属证书 (CN=agent-{uuid})                │
│  │                 │ ──► 返回: CA证书 + Agent证书 + Agent私钥                  │
│  └────────┬────────┘                                                         │
│           │                                                                  │
│           ▼                                                                  │
│  4. Agent 保存证书                                                            │
│  ┌─────────────┐                                                             │
│  │   Agent     │ ──► 保存证书到本地                                           │
│  │             │ ──► 后续通信使用 mTLS                                        │
│  └─────────────┘                                                             │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 详细设计

### 1. Join Token 设计

参考 `kubeadm join` 的 Token 设计：

```
格式: <token-id>.<token-secret>
示例: abcdef.0123456789abcdef

- token-id:     6 字符，用于 Token 查找
- token-secret: 16 字符，用于 Token 验证
```

#### Token 结构

```go
// internal/squ-apiserver/model/join_token.go

type JoinToken struct {
    ID        string    `json:"id"`         // Token ID (6字符)
    Secret    string    `json:"secret"`     // Token Secret (16字符，存储时加密)
    CreatedAt time.Time `json:"created_at"` // 创建时间
    ExpiresAt time.Time `json:"expires_at"` // 过期时间
    UsedAt    *time.Time `json:"used_at"`   // 使用时间（nil 表示未使用）
    UsedBy    string    `json:"used_by"`    // 使用该 Token 的 Agent UUID
}

// Token 完整值（用于展示给用户）
func (t *JoinToken) FullToken() string {
    return fmt.Sprintf("%s.%s", t.ID, t.Secret)
}
```

#### Token 存储

```go
// internal/squ-apiserver/repository/join_token.go

type JoinTokenRepository interface {
    Create(token *model.JoinToken) error
    GetByID(id string) (*model.JoinToken, error)
    MarkUsed(id string, agentUUID string) error
    DeleteExpired() error
}
```

### 2. squctl 命令设计

#### 生成 Token

```bash
# 创建默认 24 小时有效的 Token
squctl token create

# 创建指定有效期的 Token
squctl token create --ttl 2h
squctl token create --ttl 168h  # 7 天

# 输出示例
# Join token: abcdef.0123456789abcdef
# Expires: 2024-01-02 15:04:05 UTC
```

#### 列出 Token

```bash
squctl token list

# 输出示例
# ID       CREATED              EXPIRES              USED BY
# abcdef   2024-01-01 10:00:00  2024-01-02 10:00:00  <none>
# fedcba   2024-01-01 09:00:00  2024-01-01 10:00:00  agent-xxx
```

#### 删除 Token

```bash
squctl token delete <token-id>
```

#### 生成 Join 命令

```bash
squctl token create --print-join-command

# 输出示例
# squirrel-agent join --apiserver https://apiserver.example.com:10700 --token abcdef.0123456789abcdef
```

### 3. APIServer 端改造

#### 新增路由

```go
// internal/squ-apiserver/router/agent.go

func Agent(group *gin.RouterGroup, conf *config.Config, db database.DB) {
    // Agent 加入接口（使用 Bootstrap Token 或临时证书）
    group.POST("/agent/join", agent.JoinHandler(service))
    
    // 以下接口需要 mTLS 认证
    group.Use(mtls.Middleware(conf.MTLS))
    group.GET("/agent/certs/rotate", agent.RotateCertsHandler(service))
}
```

#### Join Handler

```go
// internal/squ-apiserver/handler/agent/service.go

type JoinRequest struct {
    Token     string `json:"token" binding:"required"` // 完整 Token
    Hostname  string `json:"hostname" binding:"required"`
    IPAddress string `json:"ip_address"`
    AgentPort int    `json:"agent_port" binding:"required"`
}

type JoinResponse struct {
    ServerID    uint   `json:"server_id"`
    UUID        string `json:"uuid"`
    CACert      []byte `json:"ca_cert"`       // CA 证书 (PEM)
    AgentCert   []byte `json:"agent_cert"`    // Agent 证书 (PEM)
    AgentKey    []byte `json:"agent_key"`     // Agent 私钥 (PEM)
    APIServerCN string `json:"apiserver_cn"`  // APIServer 证书 CN（用于验证）
}

func (s *Service) Join(request JoinRequest) response.Response {
    // 1. 解析并验证 Token
    tokenID, tokenSecret := parseToken(request.Token)
    token, err := s.TokenRepo.GetByID(tokenID)
    if err != nil {
        return response.Error(ErrTokenNotFound)
    }
    
    // 2. 检查 Token 是否过期
    if time.Now().After(token.ExpiresAt) {
        return response.Error(ErrTokenExpired)
    }
    
    // 3. 检查 Token 是否已使用（单次使用）
    if token.UsedAt != nil {
        return response.Error(ErrTokenAlreadyUsed)
    }
    
    // 4. 验证 Token Secret
    if !verifyTokenSecret(token.Secret, tokenSecret) {
        return response.Error(ErrTokenInvalid)
    }
    
    // 5. 生成 Agent UUID
    agentUUID := uuid.New().String()
    
    // 6. 生成 Agent 专属证书
    certs, err := s.CertGenerator.GenerateAgentCert(agentUUID)
    if err != nil {
        return response.Error(ErrCertGenerationFailed)
    }
    
    // 7. 创建服务器记录
    server := &model.Server{
        UUID:      agentUUID,
        Hostname:  request.Hostname,
        IpAddress: request.IPAddress,
        AgentPort: request.AgentPort,
        Status:    model.ServerStatusOnline,
    }
    if err := s.ServerRepo.Add(server); err != nil {
        return response.Error(ErrServerCreateFailed)
    }
    
    // 8. 标记 Token 已使用
    if err := s.TokenRepo.MarkUsed(tokenID, agentUUID); err != nil {
        zap.L().Error("failed to mark token as used", zap.Error(err))
    }
    
    // 9. 返回证书
    return response.Success(JoinResponse{
        ServerID:    server.ID,
        UUID:        agentUUID,
        CACert:      certs.CACert,
        AgentCert:   certs.AgentCert,
        AgentKey:    certs.AgentKey,
        APIServerCN: "squirrel-apiserver",
    })
}
```

### 4. Agent 端改造

#### 配置文件

```yaml
# config/agent.yaml
apiserver:
  # 方式一：使用 Token 加入（首次启动）
  join:
    enabled: false
    server: "https://apiserver.example.com:10700"
    token: ""    # Join Token，格式: <id>.<secret>
  
  # 方式二：使用证书通信（已注册后）
  mtls:
    enabled: true
    caFile:   "/etc/squirrel/certs/ca.crt"
    certFile: "/etc/squirrel/certs/agent.crt"
    keyFile:  "/etc/squirrel/certs/agent.key"
    serverCN: "squirrel-apiserver"  # 服务端证书 CN，用于验证
```

#### Join 流程

```go
// internal/squ-agent/server/join.go

func (s *Server) joinCluster() error {
    // 1. 检查是否已有证书
    if s.hasValidCerts() {
        zap.L().Info("Agent already has valid certificates, skipping join")
        return nil
    }
    
    // 2. 检查 Join 配置
    if !s.Config.Apiserver.Join.Enabled || s.Config.Apiserver.Join.Token == "" {
        zap.L().Info("Join not configured, waiting for manual registration")
        return nil
    }
    
    // 3. 收集主机信息
    hostInfo, err := s.collectHostInfo()
    if err != nil {
        return fmt.Errorf("failed to collect host info: %w", err)
    }
    
    // 4. 发起 Join 请求
    req := JoinRequest{
        Token:     s.Config.Apiserver.Join.Token,
        Hostname:  hostInfo.Hostname,
        IPAddress: hostInfo.IPAddress,
        AgentPort: s.Config.Server.Port,
    }
    
    resp, err := s.callJoinAPI(req)
    if err != nil {
        return fmt.Errorf("join request failed: %w", err)
    }
    
    // 5. 保存证书
    certDir := "/etc/squirrel/certs"
    if err := s.saveCerts(certDir, resp); err != nil {
        return fmt.Errorf("failed to save certificates: %w", err)
    }
    
    // 6. 保存 Server ID 到本地
    if err := s.storeServerID(resp.ServerID, resp.UUID); err != nil {
        return fmt.Errorf("failed to store server info: %w", err)
    }
    
    // 7. 更新配置，启用 mTLS
    s.Config.Apiserver.MTLS.Enabled = true
    s.Config.Apiserver.MTLS.CAFile = filepath.Join(certDir, "ca.crt")
    s.Config.Apiserver.MTLS.CertFile = filepath.Join(certDir, "agent.crt")
    s.Config.Apiserver.MTLS.KeyFile = filepath.Join(certDir, "agent.key")
    
    zap.L().Info("Agent joined cluster successfully",
        zap.String("uuid", resp.UUID),
        zap.Uint("server_id", resp.ServerID))
    
    return nil
}

func (s *Server) saveCerts(dir string, resp JoinResponse) error {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    // 保存 CA 证书
    if err := os.WriteFile(filepath.Join(dir, "ca.crt"), resp.CACert, 0644); err != nil {
        return err
    }
    
    // 保存 Agent 证书
    if err := os.WriteFile(filepath.Join(dir, "agent.crt"), resp.AgentCert, 0644); err != nil {
        return err
    }
    
    // 保存 Agent 私钥（权限更严格）
    if err := os.WriteFile(filepath.Join(dir, "agent.key"), resp.AgentKey, 0600); err != nil {
        return err
    }
    
    return nil
}
```

### 5. 证书生成器改造

扩展现有的 `internal/squctl/certs` 模块：

```go
// internal/squctl/certs/agent.go

// GenerateAgentCert 为 Agent 生成专属证书
// CN 格式: agent-{uuid}
func (g *Generator) GenerateAgentCert(uuid string) (*AgentCerts, error) {
    // 确保 CA 可用
    if err := g.ensureCA(); err != nil {
        return nil, err
    }
    
    // 生成私钥
    key, err := rsa.GenerateKey(rand.Reader, g.keySize)
    if err != nil {
        return nil, fmt.Errorf("failed to generate agent key: %w", err)
    }
    
    // 创建证书模板
    cn := fmt.Sprintf("agent-%s", uuid)
    serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    
    template := &x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            CommonName:         cn,
            Organization:       []string{"Squirrel"},
            OrganizationalUnit: []string{"Agent"},
        },
        NotBefore:   time.Now(),
        NotAfter:    time.Now().Add(g.expiry),
        KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
    }
    
    // 使用 CA 签名
    certDER, err := x509.CreateCertificate(rand.Reader, template, g.caCert, &key.PublicKey, g.caKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create agent certificate: %w", err)
    }
    
    // 读取 CA 证书
    caCertPEM, _ := os.ReadFile(filepath.Join(g.outputDir, "ca.crt"))
    
    return &AgentCerts{
        CACert:    caCertPEM,
        AgentCert: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}),
        AgentKey:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}),
    }, nil
}

type AgentCerts struct {
    CACert    []byte
    AgentCert []byte
    AgentKey  []byte
}
```

### 6. mTLS 中间件

```go
// internal/squ-apiserver/middleware/mtls.go

func Middleware(cfg *config.MTLS) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 验证客户端证书
        cert := c.Request.TLS.PeerCertificates
        if len(cert) == 0 {
            c.AbortWithStatusJSON(403, gin.H{"error": "client certificate required"})
            return
        }
        
        // 2. 验证证书 CN
        cn := cert[0].Subject.CommonName
        allowed := false
        for _, allowedCN := range cfg.AllowedCNs {
            if cn == allowedCN || strings.HasPrefix(cn, "agent-") {
                allowed = true
                break
            }
        }
        if !allowed {
            c.AbortWithStatusJSON(403, gin.H{"error": "unauthorized client certificate"})
            return
        }
        
        // 3. 将 Agent UUID 存入上下文
        if strings.HasPrefix(cn, "agent-") {
            c.Set("agent_uuid", strings.TrimPrefix(cn, "agent-"))
        }
        
        c.Next()
    }
}
```

## 使用流程

### 管理员操作

```bash
# 1. 在控制节点生成 Join Token
squctl token create --ttl 24h --print-join-command

# 输出:
# squirrel-agent join --apiserver https://192.168.1.10:10700 --token abcdef.0123456789abcdef
```

### Agent 节点操作

```bash
# 方式一：使用命令行直接加入
squirrel-agent join --apiserver https://192.168.1.10:10700 --token abcdef.0123456789abcdef

# 方式二：配置文件方式
# 编辑 /etc/squirrel/agent.yaml
# 
# apiserver:
#   join:
#     enabled: true
#     server: "https://192.168.1.10:10700"
#     token: "abcdef.0123456789abcdef"

# 启动 Agent
squirrel-agent start
```

## 安全考虑

| 安全点 | 说明 |
|--------|------|
| **Token 有效期** | 默认 24 小时，最长建议不超过 7 天 |
| **证书归属** | 每个 Agent 获得专属证书，CN 包含 UUID，便于审计 |
| **私钥保护** | Agent 私钥权限设为 0600，仅 root 可读 |
| **证书轮换** | 支持证书轮换 API，用于更新即将过期的证书 |

## 与现有证书工具集成

现有的 `internal/squctl/certs` 模块已有证书生成能力，需要扩展：

| 现有能力 | 新增能力 |
|----------|----------|
| `GenerateCA()` | 复用，用于初始化集群 |
| `GenerateServer()` | 复用，用于 APIServer 证书 |
| `GenerateClient()` | 扩展为 `GenerateAgentCert(uuid)` |
| - | 新增 `TokenCreate()` |
| - | 新增 `TokenList()` |
| - | 新增 `TokenDelete()` |

## 实现步骤

### Phase 1: Token 机制

- [ ] 实现 `JoinToken` model 和 repository
- [ ] 实现 `squctl token create/list/delete` 命令
- [ ] 实现 APIServer `/agent/join` 接口
- [ ] 实现 Agent Join 流程

### Phase 2: 证书自动分发

- [ ] 扩展 `GenerateAgentCert(uuid)` 方法
- [ ] Agent 保存证书到本地
- [ ] mTLS 中间件实现

### Phase 3: 安全加固

- [ ] Token Secret 加密存储
- [ ] 证书轮换 API
- [ ] 证书撤销机制 (CRL/OCSP)

## 参考设计

- [kubeadm join 工作原理](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/)
- [Bootstrap Tokens](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/)
