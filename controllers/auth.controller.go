package controllers

import (
	"net/http"

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

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Name: input.Name, Email: input.Email, Password: input.Password}
	user.Associate()
	user.HashPassword()

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func LoginCheck(email, password string) (string, models.User, error) {
	var err error

	user := models.User{}

	if err = database.DB.Model(models.User{}).Preload("UserReviews").Where("email=?", email).Take(&user).Error; err != nil {
		return "", user, err
	}

	err = models.VerifyPassword(password, user.Password)

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

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
