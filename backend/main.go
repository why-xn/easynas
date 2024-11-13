package main

import (
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/server"
)

func main() {
	// Initialize Zap logger
	log.InitializeLogger()

	// Connect to the database
	err := db.Connect("easynas.db")
	if err != nil {
		log.Logger.Fatal("Failed to connect to database: ", err)
	}

	err = db.GetDb().RunMigrations()
	if err != nil {
		log.Logger.Fatal("Failed to run migrations: ", err)
	}

	// Start Http Server
	server.Start()
}

/*
func testDb() {
	// Connect to the database
	dbConn, err := db.Connect("test.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the ExampleModel
	dbConn.Client().AutoMigrate(&model.User{})
	dbConn.Client().AutoMigrate(&model.NfsShare{})
	dbConn.Client().AutoMigrate(&model.NfsSharePermission{})

	// Insert a record
	user := &model.User{Name: "Shihab", Email: "shcpp@umkc.edu", IP: "10.0.8.2"}
	if err := dbConn.Insert(user); err != nil {
		log.Fatal("Insert failed:", err)
	}
	log.Println("Record inserted:", *user)

	// Insert a record
	nfsShare := &model.NfsShare{Pool: "naspool", VolName: "mydrive"}
	if err := dbConn.Insert(nfsShare); err != nil {
		log.Fatal("Insert failed:", err)
	}
	log.Println("Record inserted:", *nfsShare)

	// Get a single record
	user, err = db.Get[model.User](dbConn, map[string]interface{}{"name": "Shihab"})
	if err != nil {
		log.Fatal("Get failed:", err)
	}
	log.Println("Record retrieved:", *user)

	nfsShare, err = db.Get[model.NfsShare](dbConn, map[string]interface{}{"pool": "naspool"})
	if err != nil {
		log.Fatal("Get failed:", err)
	}
	log.Println("Record retrieved:", *nfsShare)

	// Insert a record
	nfsSharePermission := &model.NfsSharePermission{NfsShareId: nfsShare.ID, UserId: user.ID, Permission: "rw"}
	if err := dbConn.Insert(nfsSharePermission); err != nil {
		log.Fatal("Insert failed:", err)
	}
	log.Println("Record inserted:", *nfsSharePermission)

	// Update a record
	if err := dbConn.Update(&retrieved, map[string]interface{}{"age": 30}); err != nil {
		log.Fatal("Update failed:", err)
	}
	log.Println("Record updated:", *retrieved)

	// Get a list of records
	records, err := db.GetList[model.NfsSharePermission](dbConn, map[string]interface{}{"permission": "rw"}, "NfsShare", "User")
	if err != nil {
		log.Fatal("GetList failed:", err)
	}
	log.Println("Records list:", records)

	// Delete a record
	if err := dbConn.Delete(&ExampleModel{}, map[string]interface{}{"name": "Alice"}); err != nil {
		log.Fatal("Delete failed:", err)
	}
	log.Println("Record deleted")
}*/
