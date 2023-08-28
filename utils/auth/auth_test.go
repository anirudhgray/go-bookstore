package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestVerifyPassword(t *testing.T) {
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	err := VerifyPassword(plainPassword, string(hashedPassword))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		password   string
		expectPass bool
	}{
		{"Abcd123!", true},    // Strong password
		{"abcd123!", false},   // Missing uppercase
		{"ABCD123!", false},   // Missing lowercase
		{"Abcdefg!", false},   // Missing number
		{"Abcd1234", false},   // Missing special character
		{"aA1!", false},       // Too short
		{"Abcde12!", true},    // Strong password
		{"ABCD1234", false},   // Missing lowercase and special character
		{"abcd1234", false},   // Missing uppercase and special character
		{"ABCDabcd", false},   // Missing number and special character
		{"!@#$%^&*()", false}, // Missing letters and numbers
	}

	for _, test := range tests {
		result := CheckPasswordStrength(test.password)
		if result != test.expectPass {
			t.Errorf("For password %s, expected %v but got %v", test.password, test.expectPass, result)
		}
	}
}
