package mtls

import (
	"crypto/x509"
	"net/http"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
)

// MTLSAuth 返回一个 Gin 中间件，用于验证客户端 TLS 证书
// 要求服务端已启用 TLS 且配置了 ClientAuth
func MTLSAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取 TLS 连接状态
		tls := c.Request.TLS
		if tls == nil {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrTLSRequired))
			c.Abort()
			return
		}

		// 获取客户端证书
		certs := tls.PeerCertificates
		if len(certs) == 0 {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrMissingClientCert))
			c.Abort()
			return
		}

		// 获取第一个客户端证书（通常只有一个）
		clientCert := certs[0]

		// 验证证书是否过期等（TLS 握手时已经验证，这里可以提取证书信息）
		// 将证书信息存入上下文，供后续使用
		c.Set("client_cert_subject", clientCert.Subject.CommonName)
		c.Set("client_cert_issuer", clientCert.Issuer.CommonName)
		c.Set("client_cert", clientCert)

		c.Next()
	}
}

// MTLSAuthWithVerify 返回一个 Gin 中间件，验证客户端证书并检查 Common Name
func MTLSAuthWithVerify(allowedCNs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tls := c.Request.TLS
		if tls == nil {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrTLSRequired))
			c.Abort()
			return
		}

		certs := tls.PeerCertificates
		if len(certs) == 0 {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrMissingClientCert))
			c.Abort()
			return
		}

		clientCert := certs[0]
		cn := clientCert.Subject.CommonName

		// 检查 Common Name 是否在允许列表中
		if len(allowedCNs) > 0 {
			allowed := false
			for _, allowedCN := range allowedCNs {
				if cn == allowedCN {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, response.Error(response.ErrCertCNNotAllowed))
				c.Abort()
				return
			}
		}

		c.Set("client_cert_subject", cn)
		c.Set("client_cert_issuer", clientCert.Issuer.CommonName)
		c.Set("client_cert", clientCert)

		c.Next()
	}
}

// MTLSAuthWithVerifyFunc 返回一个 Gin 中间件，使用自定义函数验证客户端证书
func MTLSAuthWithVerifyFunc(verifyFunc func(cert *x509.Certificate) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tls := c.Request.TLS
		if tls == nil {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrTLSRequired))
			c.Abort()
			return
		}

		certs := tls.PeerCertificates
		if len(certs) == 0 {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrMissingClientCert))
			c.Abort()
			return
		}

		clientCert := certs[0]

		if verifyFunc != nil && !verifyFunc(clientCert) {
			c.JSON(http.StatusForbidden, response.Error(response.ErrCertVerifyFailed))
			c.Abort()
			return
		}

		c.Set("client_cert_subject", clientCert.Subject.CommonName)
		c.Set("client_cert_issuer", clientCert.Issuer.CommonName)
		c.Set("client_cert", clientCert)

		c.Next()
	}
}
