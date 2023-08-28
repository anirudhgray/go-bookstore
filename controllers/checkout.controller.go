package controllers

import (
	"net/http"
	"strconv"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
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

	if len(cart.Books) == 0 {
		c.JSON(http.StatusTeapot, gin.H{"error": "Your cart is empty!"})
		return
	}

	// check user funds
	str := "Insufficient credits. Cart Total: " + strconv.Itoa(int(cart.CalculateTotalCartPrice())) + ". Current Credit Balance: " + strconv.Itoa(currentUser.Credits)
	if cart.CalculateTotalCartPrice() > int64(currentUser.Credits) {
		c.JSON(http.StatusBadRequest, gin.H{"error": str})
		return
	}

	// Create a new transaction
	transaction := models.Transaction{
		UserID:         currentUser.ID,
		Amount:         cart.CalculateTotalCartPrice(),
		Status:         "completed",
		PaymentMethod:  "UPI",
		CreditPurchase: false,
	}

	currentUser.Credits -= int(cart.CalculateTotalCartPrice())

	for _, book := range cart.Books {
		transaction.Books = append(transaction.Books, book)
		currentUser.UserLibrary.Books = append(currentUser.UserLibrary.Books, book)
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		logger.Errorf("DB: Error creating transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	if err := database.DB.Save(&currentUser).Error; err != nil {
		logger.Errorf("DB: Error saving book to user library after transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to library"})
		return
	}

	// Clear the user's cart
	if err := database.DB.Model(&currentUser.ShoppingCart).Association("Books").Replace([]models.Book{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear shopping cart"})
		return
	}

	str = "Checkout successful. Remaining credit balance: " + strconv.Itoa(currentUser.Credits)
	c.JSON(http.StatusOK, gin.H{"message": str})
	logger.Infof("%d Credits Used on Catalog")
}

type BuyCreditsInput struct {
	Credits int `json:"credits" binding:"required"`
}

func BuyCredits(c *gin.Context) {
	currentUser := c.MustGet("user").(*models.User)

	var input BuyCreditsInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser.Credits += input.Credits

	transaction := models.Transaction{
		UserID:         currentUser.ID,
		Amount:         int64(input.Credits),
		Status:         "completed",
		PaymentMethod:  "UPI",
		CreditPurchase: true,
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		logger.Errorf("DB: Error creating credit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction."})
		return
	}

	if err := database.DB.Save(&currentUser).Error; err != nil {
		logger.Errorf("DB: Error adding credits to user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add credits."})
		return
	}

	str := "Credits added! Current balance: " + strconv.Itoa(currentUser.Credits)
	c.JSON(http.StatusOK, gin.H{"message": str})
	logger.Infof("%d Credits Purchased", input.Credits)
}
