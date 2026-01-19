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

func (j *JWT) GenToken(username string, expireDuration time.Duration) (string, error) {

	claims := CustomClaims{
		username,
		jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(expireDuration)),
			Issuer:    "Hank",
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {

	token, err := jwtgo.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtgo.Token) (i any, err error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func IsTokenExpired(tokenString string, signingKey string) (bool, error) {
	secretKey := []byte(signingKey)

	token, err := jwtgo.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtgo.Token) (any, error) {

		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return true, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.ExpiresAt.Time.Before(time.Now()), nil
	}

	return true, fmt.Errorf("invalid token")
}
