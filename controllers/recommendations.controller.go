package controllers

import (
	"net/http"
	"strconv"

	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils/recommender"
	"github.com/gin-gonic/gin"
)

func GenerateRecommendations(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(*models.User)
	userID := user.ID

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	likes, dislikes, err := recommender.GetUserLikesDislikes(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	otherUsers, err := recommender.GetUsersWithSimilarInteractions(likes, dislikes, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	similarities, err := recommender.CalculateSimilaritiesWithOtherUsers(userID, otherUsers, likes, dislikes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	unreviewed, err := recommender.GetUnreviewedBooks(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	recprobs := recommender.CalculateRecommendationProbabilities(userID, unreviewed, similarities)

	recommendedBooks := recommender.GetRecommendedBooksSortedAndPaginated(recprobs, page, 20)

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendedBooks})
}
