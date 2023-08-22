package controllers

import (
	"net/http"
	"unicode"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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

func Register(c *gin.Context) {
	// TODO email confirmation
	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !CheckPasswordStrength(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password not strong enough."})
		return
	}

	user := models.User{Name: input.Name, Email: input.Email, Password: input.Password}
	user.Associate()
	user.HashPassword()

	if err := database.DB.Create(&user).Error; err != nil {
		utils.SendRegistrationMail("Account Alert", "Someone attempted to create an account using your email. If this was you, try applying for password reset in case you have lost access to your account.", user.Email, user.ID, user.Name, false)
		c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// we lie!
		return
	}

	utils.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email, password string) (string, models.User, error) {
	var err error

	user := models.User{}

	if err = database.DB.Model(models.User{}).Preload("UserReviews").Where("email=?", email).Take(&user).Error; err != nil {
		return "", user, err
	}

	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", user, err
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return "", user, err
	}

	return token, user, nil

}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Email: input.Email, Password: input.Password}

	token, user, err := LoginCheck(user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email or password is not correct"})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email first."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
