package controllers

import (
	"fmt"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

func AddBookToCart(c *gin.Context) {
	// TODO protect against adding purchases books to cart/ buying already purchased books
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	bookID := c.Param("bookID")

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	var cart models.ShoppingCart
	if err := database.DB.Model(&cart).Preload("Books").First(&cart, "user_id = ?", currentUser.ID).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	for _, cartBook := range cart.Books {
		if cartBook.ID == book.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book is already in the cart"})
			return
		}
	}

	cart.Books = append(cart.Books, &book)

	// Save the updated shopping cart to the database
	// GORM auto saves associations on creation
	if err := database.DB.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book added to cart successfully", "cart": cart})
}

func GetCart(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var cart models.ShoppingCart
	if err := database.DB.Model(&cart).Preload("Books").First(&cart, "user_id = ?", currentUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func RemoveFromCart(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	bookID := c.Param("bookID")

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	if err := database.DB.Model(&currentUser).Preload("ShoppingCart.Books").First(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	// for deletion of many to many association, need to use association mode. Doesn't do it automatically on
	if err := database.DB.Model(&currentUser.ShoppingCart).Association("Books").Delete(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove book from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book removed from cart successfully", "cart": currentUser.ShoppingCart})
}
