package controllers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	// TODO simulate failed transaction
	currentUser := c.MustGet("user").(*models.User)

	var cart models.ShoppingCart
	if err := database.DB.Preload("Books").First(&cart, "user_id = ?", currentUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	// Create a new transaction
	transaction := models.Transaction{
		UserID:        currentUser.ID,
		Amount:        cart.CalculateTotalCartPrice(),
		Status:        "completed",
		PaymentMethod: "UPI",
	}

	for _, book := range cart.Books {
		transaction.Books = append(transaction.Books, book)
		currentUser.UserLibrary.Books = append(currentUser.UserLibrary.Books, book)
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	if err := database.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to library"})
		return
	}

	// Clear the user's cart
	if err := database.DB.Model(&currentUser.ShoppingCart).Association("Books").Replace([]models.Book{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear shopping cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
}

func GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("user").(*models.User)

	var transactions []models.Transaction
	if err := database.DB.Preload("Books").Find(&transactions, "user_id = ?", currentUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
