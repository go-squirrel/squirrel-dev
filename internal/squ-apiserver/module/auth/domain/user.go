package domain

import "context"

type CredentialVerifier interface {
	Verify(context.Context, string, string) bool
}

type TokenGenerator interface {
	Generate(string) (string, error)
}
