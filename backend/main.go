package main

import (
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/server"
)

func main() {
	// Initialize Zap logger
	log.InitializeLogger()

	// Initialize DB Connection
	err := db.Connect("easynas.db")
	if err != nil {
		log.Logger.Fatal("Failed to connect to database: ", err)
	}

	// Run DB Migrations
	err = db.GetDb().RunMigrations()
	if err != nil {
		log.Logger.Fatal("Failed to run migrations: ", err)
	}

	// Start Http Server
	server.Start()
}
