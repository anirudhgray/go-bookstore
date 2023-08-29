package auth

import (
	"unicode"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils/token"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordStrength checks for the following:
//   - min length 7
//   - uppercase
//   - lowercase
//   - number
//   - special character
func CheckPasswordStrength(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// VerifyPassword checks plaintext password against a hashed one.
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// LoginCheck checks validity of given email/password, and returns token in user exists and password is correct.
func LoginCheck(email, password string) (t string, u models.User, e error) {
	var err error

	user := models.User{}

	if err = database.DB.Model(models.User{}).Preload("UserReviews").Where("email=?", email).Take(&user).Error; err != nil {
		return "", user, err
	}

	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", user, err
	}

	token, err := token.GenerateToken(user)

	if err != nil {
		return "", user, err
	}

	return token, user, nil

}
