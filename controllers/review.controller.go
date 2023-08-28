package controllers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
)

type ReviewInput struct {
	BookID  uint   `json:"book_id" binding:"required"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
}

func AddReview(c *gin.Context) {
	currentUser := c.MustGet("user").(*models.User)

	var input ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	bookID := input.BookID

	// Check if the book exists
	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	// Check if the user has the book in their UserLibrary
	if !currentUser.HasBookInLibrary(bookID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must own the book to review it."})
		return
	}

	review := models.Review{
		UserID:  currentUser.ID,
		BookID:  bookID,
		Comment: input.Comment,
		Rating:  input.Rating,
	}

	// Check if the user has already reviewed this book
	var existingReview models.Review
	if err := database.DB.Where("user_id = ? AND book_id = ?", currentUser.ID, bookID).First(&existingReview).Error; err == nil {
		// Overwrite the existing review
		existingReview.Comment = review.Comment
		existingReview.Rating = review.Rating
		if err := database.DB.Save(&existingReview).Error; err != nil {
			logger.Errorf("DB: Review Update: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"review": existingReview})
		return
	}

	// Create a new review
	if err := database.DB.Create(&review).Error; err != nil {
		logger.Errorf("DB: Review Creation: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}
