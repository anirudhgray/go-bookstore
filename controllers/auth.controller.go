package controllers

import (
	"net/http"
	"time"

	_ "github.com/anirudhgray/balkan-assignment/docs"
	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils/auth"
	"github.com/anirudhgray/balkan-assignment/utils/email"
	"github.com/gin-gonic/gin"
)

// TODO refresh

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary register a new base user
// @Tags auth
// @Router /auth/register [post]
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
		logger.Errorf("Associate New User Failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating your account."})
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
		logger.Errorf("Attaching Cart and Library: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	c.JSON(http.StatusCreated, gin.H{"message": "User created. Verification email sent!"})
	logger.Infof("New User Created.")
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
	err := email.SendRegistrationMail("Account Verification.", "Please visit the following link to verify your account: ", user.Email, user.ID, user.Name, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in sending email."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent to you again."})
	logger.Infof("Verification requested again")
}

func ForgotPasswordRequest(c *gin.Context) {
	useremail := c.Query("email")

	var user models.User
	if err := database.DB.Where("email = ?", useremail).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Forgot Password mail sent."})
		return
	}

	var forgotPassword models.ForgotPassword
	if err := database.DB.Where("email = ?", user.Email).First(&forgotPassword).Error; err == nil {
		database.DB.Unscoped().Delete(&forgotPassword)
	}

	email.SendForgotPasswordMail(user.Email, user.ID, user.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Forgot Password mail sent."})
	logger.Infof("Forgot password request")
}

type ForgotPasswordInput struct {
	NewPassword string `json:"new_password"`
}

// after a forgot password request
func SetNewPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	useremail := c.Query("email")
	otp := c.Query("otp")
	var entry models.ForgotPassword
	if result := database.DB.Where("email = ?", useremail).First(&entry); result.Error != nil {
		logger.Errorf("Error while verifying: %v", result.Error.Error())
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid deletion verification. Account NOT deleted."})
		return
	}

	if entry.ValidTill.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "token expired, please request forgot password again."})
		return
	}

	if entry.OTP == otp {
		var user models.User
		database.DB.Where("email = ?", useremail).First(&user)

		user.Password = input.NewPassword
		user.HashPassword()

		if err := database.DB.Save(&user).Error; err != nil {
			logger.Errorf("Save user after forgot and new: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		email.GenericSendMail("Password Reset", "Password for your account was reset recently.", user.Email, user.Name)

		// database.DB.Where("email = ?", user.Email).Delete(&models.ForgotPassword{})
		c.JSON(http.StatusOK, gin.H{"message": "Password set successfully. Please proceed to login."})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"message": "Invalid verification. Password not updated."})
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
	logger.Infof("User login")
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
		logger.Errorf("Update Password failed: " + err.Error())
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
		logger.Errorf("Error while verifying.")
	}
	if entry.OTP == otp {
		var user models.User
		database.DB.Where("email = ?", email).First(&user)
		user.Verified = true
		if result := database.DB.Save(&user); result.Error != nil {
			logger.Errorf("Verification: " + result.Error.Error())
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
			logger.Errorf("Delete Acc: Failed to clear library: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear."})
			return
		}
		if err := database.DB.Model(&user.ShoppingCart).Association("Books").Clear(); err != nil {
			logger.Errorf("Delete Acc: Failed to clear cart: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear."})
			return
		}
		// Delete the user account along with associated data
		if err := database.DB.Unscoped().Select("ShoppingCart", "UserLibrary", "UserReviews").Delete(&user).Error; err != nil {
			logger.Errorf("Delete Acc: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
			return
		}
		// database.DB.Where("email = ?", user.Email).Delete(&models.DeletionConfirmation{})
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully."})
		logger.Infof("Account deleted")
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"message": "Invalid deletion verification. Account NOT deleted."})
}
