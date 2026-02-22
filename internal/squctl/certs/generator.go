package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Generator 证书生成器
type Generator struct {
	outputDir string
	keySize   int
	expiry    time.Duration

	// 已生成的 CA（用于签发其他证书）
	caCert *x509.Certificate
	caKey  *rsa.PrivateKey
}

// NewGenerator 创建证书生成器
func NewGenerator(outputDir string, keySize int, expiry time.Duration) *Generator {
	return &Generator{
		outputDir: outputDir,
		keySize:   keySize,
		expiry:    expiry,
	}
}

// GenerateCA 生成 CA 证书
func (g *Generator) GenerateCA(cn string, overwrite bool) error {
	// 检查文件是否存在
	if !overwrite && g.filesExist("ca.crt", "ca.key") {
		return fmt.Errorf("CA certificate already exists, use --overwrite to replace")
	}

	// 生成私钥
	key, err := rsa.GenerateKey(rand.Reader, g.keySize)
	if err != nil {
		return fmt.Errorf("failed to generate CA key: %w", err)
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
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(g.expiry),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 自签名
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	// 保存证书和私钥
	if err := g.saveCert("ca.crt", certDER); err != nil {
		return err
	}
	if err := g.saveKey("ca.key", key); err != nil {
		return err
	}

	// 缓存 CA 用于签发其他证书
	g.caCert = template
	g.caKey = key

	zap.L().Info("CA certificate saved",
		zap.String("cert", filepath.Join(g.outputDir, "ca.crt")),
		zap.String("key", filepath.Join(g.outputDir, "ca.key")),
	)

	return nil
}

// LoadCA 加载已有的 CA 证书
func (g *Generator) LoadCA() error {
	certPath := filepath.Join(g.outputDir, "ca.crt")
	keyPath := filepath.Join(g.outputDir, "ca.key")

	// 读取证书
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode CA certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// 读取私钥
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read CA key: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return fmt.Errorf("failed to decode CA key")
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA key: %w", err)
	}

	g.caCert = cert
	g.caKey = key

	return nil
}

// ensureCA 确保有 CA 可用
func (g *Generator) ensureCA() error {
	if g.caCert != nil && g.caKey != nil {
		return nil
	}
	return g.LoadCA()
}

// filesExist 检查文件是否存在
func (g *Generator) filesExist(filenames ...string) bool {
	for _, filename := range filenames {
		path := filepath.Join(g.outputDir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// saveCert 保存证书
func (g *Generator) saveCert(filename string, certDER []byte) error {
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	path := filepath.Join(g.outputDir, filename)
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return os.WriteFile(path, certPEM, 0644)
}

// saveKey 保存私钥
func (g *Generator) saveKey(filename string, key *rsa.PrivateKey) error {
	path := filepath.Join(g.outputDir, filename)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return os.WriteFile(path, keyPEM, 0600)
}
