package middleware

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/anirudhgray/balkan-assignment/utils/token"
	"github.com/gin-gonic/gin"
)

func BaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := token.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Authentication required"})
			logger.Errorf("Auth Middleware Error: %v", err)
			c.Abort()
			return
		}

		var user models.User

		err = database.DB.Preload("ShoppingCart").Preload("UserLibrary").First(&user, userID).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Server Error, Fetching User"})
			logger.Errorf("Fetching Authenticated User Error: %v", err)
			c.Abort()
			return
		}

		c.Set("user", &user)

		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate the token and get user details
		userID, err := token.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Authentication required"})
			logger.Errorf("Auth Middleware Error: %v", err)
			c.Abort()
			return
		}

		var user models.User

		err = database.DB.Preload("ShoppingCart").Preload("UserLibrary").First(&user, userID).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Server Error, Fetching User"})
			logger.Errorf("Fetching Authenticated User Error: %v", err)
			c.Abort()
			return
		}

		c.Set("user", &user)

		// Check if the user is authorized as an admin
		if user.Role != models.Admin {
			c.JSON(http.StatusForbidden, gin.H{"Forbidden": "Access denied"})
			logger.Errorf("Admin Auth Middleware Error: User is not authorized as admin")
			c.Abort()
			return
		}

		c.Next()
	}
}
