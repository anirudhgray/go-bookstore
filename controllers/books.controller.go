package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

// AttachCL is a helper for attaching cart/library/wishlist to old accounts
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
