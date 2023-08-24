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
		Name:          input.Name,
		Author:        input.Author,
		Publisher:     input.Publisher,
		Date:          input.Date,
		Price:         input.Price,
		ISBN:          input.ISBN,
		Category:      models.BookCategory(input.Category), // Convert string to enum value
		FilePath:      filePath,
		CatalogDelete: false,
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

// TODO GetAllBooks, including catalog_deleted ones.

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

// gets all shopping cart transactions done
func GetAllTransactions(c *gin.Context) {
	var transactions []models.Transaction
	if err := database.DB.Preload("Books").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all transactions for admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// delete a review. Note that due to one review per user constraint, use cannot add another review, even once deleted (because soft delete). Not sure if this needs to be patched, might be a feature.
func DeleteReview(c *gin.Context) {
	var review models.Review
	reviewID := c.Param("id")

	if err := database.DB.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if err := database.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// DeleteBook deletes a book from from catalog and carts, not from user libraries.
func DeleteBook(c *gin.Context) {
	bookID := c.Param("id")

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	book.CatalogDelete = true

	if err := database.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book has been marked as deleted"})
}

// DeleteBookHard deletes a book from user libraries as well. Do not use unless necessary.
func DeleteBookHard(c *gin.Context) {
	var book models.Book
	bookID := c.Param("id")

	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// TODO implement user banning once user deactivation implmented
