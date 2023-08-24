package migrations

import (
	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
)

func Migrate() {
	var migrationModels = []interface{}{
		&models.User{},
		&models.ShoppingCart{},
		&models.Book{},
		&models.UserLibrary{},
		&models.Review{},
		&models.VerificationEntry{},
		&models.Transaction{},
		&models.DeletionConfirmation{},
	}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
