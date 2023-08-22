package controllers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Right now, I'm doing a lot in the controllers.
// - http request handling
// - database operations
// - business logic
// should separate out db ops and http handlers.

func GetData(ctx *gin.Context) {
	var example []*models.Example
	result := database.DB.Find(&example)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}
	ctx.JSON(http.StatusOK, &example)

}

func GetSingleData(ctx *gin.Context) {
	exampleID := ctx.Param("pid")

	example := new(models.Example)
	result := database.DB.First(example, exampleID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	ctx.JSON(http.StatusOK, &example)
}

func Create(ctx *gin.Context) {
	example := new(models.Example)

	// Bind the request JSON and validate the existence of the 'data' attribute
	if err := ctx.ShouldBindJSON(&example); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Save the example to the database
	result := database.DB.Create(example)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create example"})
		return
	}

	ctx.JSON(http.StatusOK, example)
}

func Update(ctx *gin.Context) {
	// Get the example ID from the URL parameter
	exampleID := ctx.Param("pid")

	// Find the existing example by ID
	existingExample := new(models.Example)
	result := database.DB.First(existingExample, exampleID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	// Bind the JSON data to a temporary struct for partial updates
	var updateData struct {
		gorm.Model
		Data string
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	database.DB.Model(existingExample).Updates(updateData)

	ctx.JSON(http.StatusOK, existingExample)
}
