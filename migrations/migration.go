package migrations

import (
	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{&models.Example{}}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
