package controllers

import (
	"net/http"
	"strconv"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

func AddBookToCart(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	id := c.Param("bookID")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400"})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	if book.CatalogDelete {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found."})
		return
	}

	var cart models.ShoppingCart
	if err := database.DB.Model(&cart).Preload("Books").First(&cart, "user_id = ?", currentUser.ID).Error; err != nil {
		logger.Errorf("Fetch cart: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	for _, cartBook := range cart.Books {
		if cartBook.ID == book.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book is already in the cart"})
			return
		}
	}

	var library models.UserLibrary
	if err := database.DB.Model(&library).Preload("Books").First(&library, "user_id = ?", currentUser.ID).Error; err != nil {
		logger.Errorf("Fetch library: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user library"})
		return
	}

	// check if book in purchased library already
	for _, libraryBook := range library.Books {
		if libraryBook.ID == book.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book is already purchased and in user library"})
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
		logger.Errorf("Getting cart: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func RemoveFromCart(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	id := c.Param("bookID")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400"})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		logger.Errorf("Removing from cart, fetch book: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	if err := database.DB.Model(&currentUser).Preload("ShoppingCart.Books").First(&currentUser).Error; err != nil {
		logger.Errorf("Getting user, preload cart books: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping cart"})
		return
	}

	// for deletion of many to many association, need to use association mode. Doesn't do it automatically on
	if err := database.DB.Model(&currentUser.ShoppingCart).Association("Books").Delete(book); err != nil {
		logger.Errorf("Remove from cart: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove book from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book removed from cart successfully", "cart": currentUser.ShoppingCart})
}
