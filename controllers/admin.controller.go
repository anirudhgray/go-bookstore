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

// GetAllUsers gets undeleted users.
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Fetch all users from the database
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
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

	if err := database.DB.Select("Reviews").Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func BanUser(c *gin.Context) {
	userID := c.Param("userID")
	status := c.Query("status")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Ban/Unban the user
	user.Banned = status == "true"
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change ban status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User ban status changed successfully"})
}

func PromoteUserToAdmin(c *gin.Context) {
	// Get the user ID of user to promote
	userID := c.Param("userID")
	currentUser := c.MustGet("user").(*models.User)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Email == currentUser.Email {
		c.JSON(http.StatusTeapot, gin.H{"error": "You can't promote yourself!"})
		return
	}

	user.Role = models.Admin
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Specified user is now an admin"})
}

func DemoteUserToBase(c *gin.Context) {
	// Get the user ID of user to promote
	userID := c.Param("userID")
	currentUser := c.MustGet("user").(*models.User)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Email == currentUser.Email {
		c.JSON(http.StatusTeapot, gin.H{"error": "You can't demote yourself!"})
		return
	}

	user.Role = models.BaseUser
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Specified user is now a base user."})
}
