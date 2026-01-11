package hash

import (
	"fmt"
	"log"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password string
		salt     string
		expected string
	}{
		{"password123", "", "Expected hashed password"},
		{"password123", "saltsalt", "Expected hashed password with salt"},
		{"", "", "password cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Password: %s Salt: %s", tt.password, tt.salt), func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password, tt.salt)

			if err != nil && tt.password == "" {
				if err.Error() != tt.expected {
					t.Errorf("expected error %s, got %s", tt.expected, err.Error())
				}
			} else if err == nil {
				if hashedPassword == "" {
					t.Errorf("expected non-empty hashed password, got empty")
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	err = ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("password should match, but got error: %v", err)
	}

	err = ComparePassword(hashedPassword, "wrongpassword")
	fmt.Println(err)
	if err == nil {
		t.Error("password should not match, but it did")
	}
}
