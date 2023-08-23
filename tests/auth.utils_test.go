package tests

import (
	"testing"

	"github.com/anirudhgray/balkan-assignment/utils/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestVerifyPassword(t *testing.T) {
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	err := auth.VerifyPassword(plainPassword, string(hashedPassword))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
