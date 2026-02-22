package certs

import (
	"time"

	"go.uber.org/zap"
)

// DefaultExpiry 默认证书有效期（10年）
const DefaultExpiry = 87600 * time.Hour

// CertsOptions 证书生成选项
type CertsOptions struct {
	// 输出目录
	OutputDir string

	// CA 相关
	OnlyCA bool
	CACN   string

	// Server 证书相关
	OnlyServer   bool
	ServerCN     string
	ServerHosts  []string

	// Client 证书相关
	OnlyClient bool
	ClientCN   string

	// 通用参数
	Expiry    time.Duration
	KeySize   int
	Overwrite bool

	// 内部状态
	generator *Generator
}

// NewCertsOptions 创建证书选项
func NewCertsOptions() *CertsOptions {
	return &CertsOptions{
		OutputDir:   "./certs",
		CACN:        "squirrel-ca",
		ServerCN:    "squirrel-apiserver",
		ServerHosts: []string{"127.0.0.1", "localhost"},
		ClientCN:    "squirrel-agent",
		Expiry:      DefaultExpiry,
		KeySize:     2048,
		Overwrite:   false,
	}
}

// Run 执行证书生成
func (o *CertsOptions) Run() error {
	// 初始化生成器
	o.generator = NewGenerator(o.OutputDir, o.KeySize, o.Expiry)

	// 根据选项生成不同的证书
	switch {
	case o.OnlyCA:
		return o.generateCA()
	case o.OnlyServer:
		return o.generateServer()
	case o.OnlyClient:
		return o.generateClient()
	default:
		return o.generateAll()
	}
}

// generateCA 仅生成 CA 证书
func (o *CertsOptions) generateCA() error {
	zap.L().Info("generating CA certificate",
		zap.String("cn", o.CACN),
		zap.String("output_dir", o.OutputDir),
	)

	if err := o.generator.GenerateCA(o.CACN, o.Overwrite); err != nil {
		zap.L().Error("failed to generate CA certificate", zap.Error(err))
		return err
	}

	zap.L().Info("CA certificate generated successfully")
	return nil
}

// generateServer 仅生成服务端证书
func (o *CertsOptions) generateServer() error {
	zap.L().Info("generating server certificate",
		zap.String("cn", o.ServerCN),
		zap.Strings("hosts", o.ServerHosts),
	)

	if err := o.generator.GenerateServer(o.ServerCN, o.ServerHosts, o.Overwrite); err != nil {
		zap.L().Error("failed to generate server certificate", zap.Error(err))
		return err
	}

	zap.L().Info("server certificate generated successfully")
	return nil
}

// generateClient 仅生成客户端证书
func (o *CertsOptions) generateClient() error {
	zap.L().Info("generating client certificate",
		zap.String("cn", o.ClientCN),
	)

	if err := o.generator.GenerateClient(o.ClientCN, o.Overwrite); err != nil {
		zap.L().Error("failed to generate client certificate", zap.Error(err))
		return err
	}

	zap.L().Info("client certificate generated successfully")
	return nil
}

// generateAll 生成所有证书
func (o *CertsOptions) generateAll() error {
	zap.L().Info("generating all certificates",
		zap.String("output_dir", o.OutputDir),
	)

	// 1. 生成 CA
	if err := o.generateCA(); err != nil {
		return err
	}

	// 2. 生成服务端证书
	if err := o.generateServer(); err != nil {
		return err
	}

	// 3. 生成客户端证书
	if err := o.generateClient(); err != nil {
		return err
	}

	zap.L().Info("all certificates generated successfully")
	return nil
}
