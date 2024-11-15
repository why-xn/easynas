package db

import (
	"github.com/whyxn/easynas/backend/pkg/db/model"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	golog "log"
	"os"
	"time"
)

const RecordNotFound = "record not found"

// Database struct to hold the connection instance
type Database struct {
	client *gorm.DB
}

var db = Database{}

// Connect initializes a connection to the SQLite database
func Connect(dbPath string) error {
	var err error

	newLogger := logger.New(
		golog.New(os.Stdout, "\r\n", golog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)

	db.client, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	log.Logger.Info("Database connection established")
	return nil
}

func GetDb() *Database {
	return &db
}

// Client returns the current database client
func (db *Database) Client() *gorm.DB {
	return db.client
}

func (db *Database) RunMigrations() error {
	log.Logger.Info("Running Database Migrations...")

	var err error
	// Auto-migrate the ExampleModel
	err = db.Client().AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	err = db.Client().AutoMigrate(&model.NfsShare{})
	if err != nil {
		return err
	}

	err = db.Client().AutoMigrate(&model.NfsSharePermission{})
	if err != nil {
		return err
	}

	// Create Initial Admin User
	// Check if admin user already exists in the DB
	_, err = Get[model.User](db, map[string]interface{}{"email": "admin@easy.nas"})
	if err != nil && err.Error() == RecordNotFound {
		// Create and insert initial admin user in DB
		hashPassword, err := util.HashPassword("admin")
		if err != nil {
			log.Logger.Fatalw("Failed to hash admin password", "err", err.Error())
		}
		user := &model.User{
			Name:        "Admin",
			Email:       "admin@easy.nas",
			Password:    hashPassword,
			NasClientIP: "10.0.0.1",
			Role:        model.RoleAdmin,
		}
		if err = db.Insert(user); err != nil {
			log.Logger.Fatalw("Failed to create initial admin user", "err", err.Error())
		}
	}

	return nil
}

// Insert inserts a record into the specified table
func (db *Database) Insert(record interface{}) error {
	result := db.client.Create(record)
	return result.Error
}

// Get retrieves a single record based on the given conditions and returns it as an object of the specified type
func Get[T any](db *Database, conditions map[string]interface{}, preloads ...string) (*T, error) {
	var result T
	var tx = db.client

	if len(preloads) > 0 {
		for _, preload := range preloads {
			tx = tx.Preload(preload)
		}
	}

	if err := tx.Where(conditions).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// GetList retrieves multiple records based on the given conditions and returns them as a slice of objects of the specified type
func GetList[T any](db *Database, conditions map[string]interface{}, preloads ...string) ([]T, error) {
	var results []T
	var tx = db.client

	if len(preloads) > 0 {
		for _, preload := range preloads {
			tx = tx.Preload(preload)
		}
	}

	if err := tx.Where(conditions).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Update modifies fields in a record matching the specified conditions
func (db *Database) Update(record interface{}, updates map[string]interface{}) error {
	result := db.client.Model(record).Updates(updates)
	return result.Error
}

// Delete removes a record matching the specified conditions from the table
func (db *Database) Delete(record interface{}, conditions map[string]interface{}) error {
	result := db.client.Where(conditions).Delete(record)
	return result.Error
}
