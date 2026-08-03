package infra

import (
	"context"

	"gorm.io/gorm"

	"squirrel-dev/pkg/hash"
)

type Verifier struct{ db *gorm.DB }

func NewVerifier(db *gorm.DB) *Verifier { return &Verifier{db: db} }

func (v *Verifier) Verify(ctx context.Context, username, password string) bool {
	var user userModel
	if err := v.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return hash.ComparePassword(user.Password, password) == nil
}
