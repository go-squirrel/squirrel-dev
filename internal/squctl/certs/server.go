package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// GenerateServer 生成服务端证书
func (g *Generator) GenerateServer(cn string, hosts []string, overwrite bool) error {
	// 确保 CA 可用
	if err := g.ensureCA(); err != nil {
		return fmt.Errorf("CA certificate not found, generate CA first: %w", err)
	}

	// 检查文件是否存在
	if !overwrite && g.filesExist("server.crt", "server.key") {
		return fmt.Errorf("server certificate already exists, use --overwrite to replace")
	}

	// 生成私钥
	key, err := rsa.GenerateKey(rand.Reader, g.keySize)
	if err != nil {
		return fmt.Errorf("failed to generate server key: %w", err)
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
			OrganizationalUnit: []string{"APIServer"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(g.expiry),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	// 解析 hosts 为 IP 和 DNS
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	// 使用 CA 签名
	certDER, err := x509.CreateCertificate(rand.Reader, template, g.caCert, &key.PublicKey, g.caKey)
	if err != nil {
		return fmt.Errorf("failed to create server certificate: %w", err)
	}

	// 保存证书和私钥
	if err := g.saveCert("server.crt", certDER); err != nil {
		return err
	}
	if err := g.saveKey("server.key", key); err != nil {
		return err
	}

	zap.L().Info("server certificate saved",
		zap.String("cert", filepath.Join(g.outputDir, "server.crt")),
		zap.String("key", filepath.Join(g.outputDir, "server.key")),
		zap.Strings("hosts", hosts),
	)

	return nil
}
