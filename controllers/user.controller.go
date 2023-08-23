package controllers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

func GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("user").(*models.User)

	var transactions []models.Transaction
	if err := database.DB.Preload("Books").Find(&transactions, "user_id = ?", currentUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func GetUserLibrary(c *gin.Context) {
	currentUser := c.MustGet("user").(*models.User)

	var userLibrary []models.UserLibrary
	if err := database.DB.Preload("Books").Find(&userLibrary, "user_id = ?", currentUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"library": userLibrary})
}
