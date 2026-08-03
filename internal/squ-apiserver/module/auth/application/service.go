package application

import (
	"context"

	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/auth/domain"
)

type Service struct {
	verifier domain.CredentialVerifier
	tokens   domain.TokenGenerator
}

func NewService(verifier domain.CredentialVerifier, tokens domain.TokenGenerator) *Service {
	return &Service{verifier: verifier, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	if !s.verifier.Verify(ctx, username, password) {
		zap.L().Warn("invalid login credentials", zap.String("username", username))
		return "", ErrInvalidCredentials
	}
	token, err := s.tokens.Generate(username)
	if err != nil {
		zap.L().Error("failed to generate token", zap.String("username", username), zap.Error(err))
		return "", ErrTokenGeneration
	}
	return token, nil
}
