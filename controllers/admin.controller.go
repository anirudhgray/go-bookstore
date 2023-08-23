package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

type CreateBookInput struct {
	Name      string    `json:"name" binding:"required"`
	Author    string    `json:"author" binding:"required"`
	Publisher string    `json:"publisher" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Price     int64     `json:"price" binding:"required"` // TODO add localisation/internationalisation for price
	ISBN      string    `json:"isbn" binding:"required"`
	Category  string    `json:"category" binding:"required"`
}

// CreateBook creates a new book. Admins only.
func CreateBook(c *gin.Context) {
	var input CreateBookInput

	// Validate request data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new Book object
	book := models.Book{
		Name:      input.Name,
		Author:    input.Author,
		Publisher: input.Publisher,
		Date:      input.Date,
		Price:     input.Price,
		ISBN:      input.ISBN,
		Category:  models.BookCategory(input.Category), // Convert string to enum value
		FilePath:  filePath,
	}

	// Save the book to the database
	if err := database.DB.Create(&book).Error; err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully", "data": book})
}

type EditBookInput struct {
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	Publisher string    `json:"publisher"`
	Date      time.Time `json:"date"` // Assuming a string format for simplicity
	Price     int64     `json:"price"`
	ISBN      string    `json:"isbn"`
	Category  string    `json:"category"`
}

func EditBook(c *gin.Context) {
	bookID := c.Param("id") // Extract the book ID from the URL parameter

	var input EditBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := database.DB.Model(&book).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "book": book})
}

// TODO add user functions such as deleting users, seeing statistics for how many users, how many deactivated, how many comments etc.
