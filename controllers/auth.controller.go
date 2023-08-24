package controllers

import (
	"net/http"
	"time"

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
	if err := user.Associate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	user.HashPassword()

	if err := database.DB.Create(&user).Error; err != nil {
		email.SendRegistrationMail("Account Alert", "Someone attempted to create an account using your email. If this was you, try applying for password reset in case you have lost access to your account.", user.Email, user.ID, user.Name, false)
		c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// we lie!
		return
	}

	if err := user.AttachCartAndLibrary(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
}

func RequestVerificationAgain(c *gin.Context) {
	useremail := c.Query("email")

	var user models.User
	if err := database.DB.Where("email = ?", useremail).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Verification email sent."})
		return
	}

	if user.Verified {
		c.JSON(http.StatusOK, gin.H{"message": "Verification email sent."})
		return
	}

	// Check if a deletion confirmation record already exists for the user's email
	var verificationEntry models.VerificationEntry
	if err := database.DB.Where("email = ?", user.Email).First(&verificationEntry).Error; err == nil {
		database.DB.Unscoped().Delete(&verificationEntry)
	}

	// Send deletion email
	email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent."})
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

	if user.Banned {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are banned! Contact admin if you think this is a mistake."})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email first."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

type ResetPasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordController handles the reset password by logged in user
func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	if err := auth.VerifyPassword(input.OldPassword, currentUser.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
		email.GenericSendMail("Password Reset Attempt", "Somebody attempted to change your password on Bookstore. Secure your account if this was not you.", currentUser.Email, currentUser.Name)
		return
	}

	currentUser.Password = input.NewPassword
	currentUser.HashPassword()

	if err := database.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	email.GenericSendMail("Password Reset Successfully", "Your password for Bookstore was changed. Secure your account if this was not you.", currentUser.Email, currentUser.Name)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
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
		var verificationEntry models.VerificationEntry
		if err := database.DB.Where("email = ?", email).First(&verificationEntry).Error; err == nil {
			database.DB.Unscoped().Delete(&verificationEntry)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Verified! You can now log in."})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification."})
}

func RequestDeletion(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	// Check if a deletion confirmation record already exists for the user's email
	var existingConfirmation models.DeletionConfirmation
	if err := database.DB.Where("email = ?", currentUser.Email).First(&existingConfirmation).Error; err == nil {
		database.DB.Unscoped().Delete(&existingConfirmation)
	}

	// Send deletion email
	email.SendDeletionMail(currentUser.Email, currentUser.ID, currentUser.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Deletion request submitted"})
}

func DeleteAccount(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")
	var entry models.DeletionConfirmation
	if result := database.DB.Where("email = ?", email).First(&entry); result.Error != nil {
		logger.Errorf("Error while verifying: %v", result.Error)
	}

	if entry.ValidTill.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "token expired, please request again."})
		return
	}

	if entry.OTP == otp {
		var user models.User
		database.DB.Where("email = ?", email).Preload("UserLibrary").Preload("ShoppingCart").First(&user)

		if err := database.DB.Model(&user.UserLibrary).Association("Books").Clear(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear."})
			return
		}
		if err := database.DB.Model(&user.ShoppingCart).Association("Books").Clear(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear."})
			return
		}
		// Delete the user account along with associated data
		if err := database.DB.Unscoped().Select("ShoppingCart", "UserLibrary", "UserReviews").Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
			return
		}
		database.DB.Where("email = ?", user.Email).Delete(&models.DeletionConfirmation{})
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully."})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"message": "Invalid deletion verification. Account NOT deleted."})
}
