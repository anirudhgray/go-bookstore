package main

import (
	"time"

	"github.com/anirudhgray/balkan-assignment/config"
	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/migrations"
	"github.com/anirudhgray/balkan-assignment/routers"
	"github.com/spf13/viper"
)

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	migrations.Migrate()

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}

// TODO add cron job or route to clean up deleted books, verification entries, reviews, etc.
// TODO initially just soft delete user. Do current hard delete routine after 24 hours.
// TODO You trusted all proxies, this is NOT safe. We recommend you to set a value.
