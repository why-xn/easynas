package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

// Database struct to hold the connection instance
type Database struct {
	client *gorm.DB
}

// Connect initializes a connection to the SQLite database
func Connect(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connection established.")
	return &Database{client: db}, nil
}

// Client returns the current database client
func (db *Database) Client() *gorm.DB {
	return db.client
}

// Insert inserts a record into the specified table
func (db *Database) Insert(record interface{}) error {
	result := db.client.Create(record)
	return result.Error
}

// Get retrieves a single record based on the given conditions and returns it as an object of the specified type
func Get[T any](db *Database, conditions map[string]interface{}) (*T, error) {
	var result T
	if err := db.client.Where(conditions).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// GetList retrieves multiple records based on the given conditions and returns them as a slice of objects of the specified type
func GetList[T any](db *Database, conditions map[string]interface{}) ([]T, error) {
	var results []T
	if err := db.client.Where(conditions).Find(&results).Error; err != nil {
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
