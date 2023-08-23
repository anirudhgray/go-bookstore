package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

func AttachCL(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)
	currentUser.AttachCartAndLibrary()
	c.JSON(http.StatusOK, gin.H{"message": "cart and library attached"})
}

func GetBooks(c *gin.Context) {
	var books []models.Book

	query := database.DB.Model(&models.Book{})

	// Pagination
	page := c.DefaultQuery("page", "1") // Get the requested page (default to 1)
	pageSize := 20
	i, err := strconv.Atoi(page) // Number of books per page
	if err != nil {
		logger.Errorf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error."})
		return
	}
	offset := (i - 1) * pageSize

	// Search by name, author, or category
	searchQuery := c.DefaultQuery("search", "")
	if searchQuery != "" {
		searchQuery = strings.ToLower(searchQuery) // Convert search query to lowercase
		query = query.Where("LOWER(name) LIKE ? OR LOWER(author) LIKE ? OR LOWER(category) LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Filter by category
	categoryFilter := c.DefaultQuery("category", "")
	if categoryFilter != "" {
		query = query.Where("category = ?", categoryFilter)
	}

	// Sort by price
	sortByPrice := c.DefaultQuery("sortByPrice", "asc")
	if sortByPrice == "asc" {
		query = query.Order("price ASC")
	} else if sortByPrice == "desc" {
		query = query.Order("price DESC")
	}

	// Apply pagination
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

func AddBookToCart(c *gin.Context) {
	// Get the logged-in user from the context
	user, _ := c.Get("user")
	currentUser := user.(*models.User) // Assert user type

	// Get the book ID from the URL parameter
	bookID := c.Param("bookID")

	// Fetch the book from the database
	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	// Fetch the user's shopping cart
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

	// Add the book to the shopping cart's books slice
	cart.Books = append(cart.Books, &book)

	// Save the updated shopping cart to the database
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
