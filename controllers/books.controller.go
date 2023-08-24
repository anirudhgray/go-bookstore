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

type BookWithAvgRating struct {
	models.Book
	AvgRating float64
}

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
		searchQuery = strings.ToLower(searchQuery)
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

	// remove catalog deleted books
	query = query.Where("catalog_delete = ?", false)

	if err := query.Preload("Reviews").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	booksWithAvgRating := make([]BookWithAvgRating, 0)
	for _, book := range books {
		logger.Infof("%v %v", book.Name, book.CalculateAvgRating())

		bookWithAvgRating := BookWithAvgRating{
			Book:      book,
			AvgRating: book.CalculateAvgRating(),
		}
		booksWithAvgRating = append(booksWithAvgRating, bookWithAvgRating)
	}

	c.JSON(http.StatusOK, gin.H{"books": booksWithAvgRating})
}

func GetBook(c *gin.Context) {
	var book models.Book
	bookID := c.Param("bookID")

	if err := database.DB.Preload("Reviews").First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.CatalogDelete {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	bookWithAvgRating := BookWithAvgRating{
		Book:      book,
		AvgRating: book.CalculateAvgRating(),
	}

	// Calculate average rating for the book
	book.CalculateAvgRating()

	c.JSON(http.StatusOK, gin.H{"book": bookWithAvgRating})
}

func DownloadBook(c *gin.Context) {
	bookID := c.Param("bookID")

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	currentUser, _ := c.Get("user")

	user := currentUser.(*models.User)
	if !user.HasBookInLibrary(book.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this book"})
		return
	}

	// Send the file using c.File()
	if book.FilePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No pdf found. Contact support."})
		return
	}
	c.File(book.FilePath)
}
