package controllers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils/recommender"
	"github.com/gin-gonic/gin"
)

func GenerateRecommendations(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(*models.User)
	userID := user.ID

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

	c.JSON(http.StatusOK, gin.H{"others": otherUsers})
}
