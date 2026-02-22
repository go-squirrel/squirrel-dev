package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// GenerateClient 生成客户端证书
func (g *Generator) GenerateClient(cn string, overwrite bool) error {
	// 确保 CA 可用
	if err := g.ensureCA(); err != nil {
		return fmt.Errorf("CA certificate not found, generate CA first: %w", err)
	}

	// 检查文件是否存在
	if !overwrite && g.filesExist("client.crt", "client.key") {
		return fmt.Errorf("client certificate already exists, use --overwrite to replace")
	}

	// 生成私钥
	key, err := rsa.GenerateKey(rand.Reader, g.keySize)
	if err != nil {
		return fmt.Errorf("failed to generate client key: %w", err)
	}

	// 创建证书模板
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"Squirrel"},
			Country:      []string{"CN"},
			Province:     []string{"Beijing"},
			Locality:     []string{"Beijing"},
			OrganizationalUnit: []string{"Agent"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(g.expiry),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// 使用 CA 签名
	certDER, err := x509.CreateCertificate(rand.Reader, template, g.caCert, &key.PublicKey, g.caKey)
	if err != nil {
		return fmt.Errorf("failed to create client certificate: %w", err)
	}

	// 保存证书和私钥
	if err := g.saveCert("client.crt", certDER); err != nil {
		return err
	}
	if err := g.saveKey("client.key", key); err != nil {
		return err
	}

	zap.L().Info("client certificate saved",
		zap.String("cert", filepath.Join(g.outputDir, "client.crt")),
		zap.String("key", filepath.Join(g.outputDir, "client.key")),
		zap.String("cn", cn),
	)

	return nil
}
