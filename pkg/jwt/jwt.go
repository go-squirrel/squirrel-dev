package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SigningKey []byte
}

func New(signingKey string) *JWT {
	return &JWT{
		[]byte(signingKey),
	}
}

type CustomClaims struct {
	// UUID        uuid.UUID
	Username string
	jwtgo.RegisteredClaims
}

// GenToken 生成JWT
func (j *JWT) GenToken(username string, expireDuration time.Duration) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		username, // 自定义字段
		jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(expireDuration)),
			Issuer:    "Hank", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析JWT
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwtgo.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtgo.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func IsTokenExpired(tokenString string, signingKey string) (bool, error) {
	secretKey := []byte(signingKey)
	// 解析 token
	token, err := jwtgo.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return true, err
	}

	// 获取 token 的声明
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 检查是否过期
		return claims.ExpiresAt.Time.Before(time.Now()), nil
	}

	return true, fmt.Errorf("invalid token")
}
