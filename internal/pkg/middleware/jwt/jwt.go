package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/pkg/jwt"
)

// JWTAuth 返回一个 Gin 中间件，用于验证 JWT Token
func JWTAuth(signingKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrMissingAuthHeader))
			c.Abort()
			return
		}

		// 提取 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrInvalidAuthFormat))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析并验证 token
		j := jwt.New(signingKey)
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error(response.ErrTokenInvalid))
			c.Abort()
			return
		}

		// 将用户名存入上下文，供后续使用
		c.Set("username", claims.Username)

		c.Next()
	}
}
