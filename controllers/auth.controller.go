package controllers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils/auth"
	"github.com/anirudhgray/balkan-assignment/utils/email"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TODO forgot password

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !auth.CheckPasswordStrength(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password not strong enough."})
		return
	}

	user := models.User{Name: input.Name, Email: input.Email, Password: input.Password}
	user.Associate()
	user.HashPassword()

	if err := database.DB.Create(&user).Error; err != nil {
		email.SendRegistrationMail("Account Alert", "Someone attempted to create an account using your email. If this was you, try applying for password reset in case you have lost access to your account.", user.Email, user.ID, user.Name, false)
		c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// we lie!
		return
	}

	email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Email: input.Email, Password: input.Password}

	token, user, err := auth.LoginCheck(user.Email, user.Password)

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

// VerifyEmail takes your email and otp sent of registration to verify a user account.
func VerifyEmail(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")
	var entry models.VerificationEntry
	if result := database.DB.Where("email = ?", email).First(&entry); result.Error != nil {
		logger.Errorf("Error while verifying: %v", result.Error)
	}
	if entry.OTP == otp {
		var user models.User
		database.DB.Where("email = ?", email).First(&user)
		user.Verified = true
		if result := database.DB.Save(&user); result.Error != nil {
			logger.Errorf("Error while verifying: %v", result.Error)
			return
		}
		database.DB.Where("email = ?", email).Delete(&models.VerificationEntry{})
		c.JSON(http.StatusOK, gin.H{"message": "Verified! You can now log in."})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification."})
}
