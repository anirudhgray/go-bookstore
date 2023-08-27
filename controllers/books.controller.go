package controllers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

type BookWithAvgRating struct {
	Book      models.SafeBook
	AvgRating float64
}

// AttachCL is a helper for attaching cart/library/wishlist to old accounts
func AttachCL(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)
	currentUser.AttachCartAndLibrary()
	c.JSON(http.StatusOK, gin.H{"message": "cart and library attached"})
}

// Paginated,
// Search by name author, category or exact ISBN
// Filter by category,
// Sort by price or newness
func GetBooks(c *gin.Context) {
	var books []models.SafeBook

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

	// Search by name, author, category, or exact ISBN
	searchQuery := c.DefaultQuery("search", "")
	if searchQuery != "" {
		searchQuery = strings.ToLower(searchQuery)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(author) LIKE ? OR LOWER(category) LIKE ? OR ISBN = ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%", searchQuery)
	}

	// Filter by category
	categoryFilter := c.DefaultQuery("category", "")
	if categoryFilter != "" {
		query = query.Where("category = ?", categoryFilter)
	}

	// Sort by price
	sortByPrice := c.DefaultQuery("sortByPrice", "")
	if sortByPrice == "asc" {
		query = query.Order("price ASC")
	} else if sortByPrice == "desc" {
		query = query.Order("price DESC")
	}

	// Sort by newly added
	sortByNewness := c.DefaultQuery("sortByNewness", "")
	if sortByNewness == "asc" {
		query = query.Order("created_at ASC")
	} else if sortByNewness == "desc" {
		query = query.Order("created_at DESC")
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
		logger.Infof("%v %v\n", book.Name, book.CalculateAvgRating())

		bookWithAvgRating := BookWithAvgRating{
			Book:      book,
			AvgRating: book.CalculateAvgRating(),
		}
		booksWithAvgRating = append(booksWithAvgRating, bookWithAvgRating)
	}

	c.JSON(http.StatusOK, gin.H{"books": booksWithAvgRating})
}

func GetBook(c *gin.Context) {
	var book models.SafeBook
	id := c.Param("bookID")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400"})
		return
	}

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
	id := c.Param("bookID")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400"})
		return
	}

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

	// fetch file from cloud storage loc
	res, err := http.Get(book.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Trouble with cloud storage."})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch remote file."})
		return
	}

	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, res.Body)

	// c.File(book.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
