package infra

import (
	"time"

	"squirrel-dev/pkg/jwt"
)

type TokenGenerator struct {
	signingKey string
	expires    time.Duration
}

func NewTokenGenerator(signingKey string, expiredMinutes int) *TokenGenerator {
	return &TokenGenerator{
		signingKey: signingKey,
		expires:    time.Duration(expiredMinutes) * time.Minute,
	}
}

func (g *TokenGenerator) Generate(username string) (string, error) {
	return jwt.New(g.signingKey).GenToken(username, g.expires)
}
