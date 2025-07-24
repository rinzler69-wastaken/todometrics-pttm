package config

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-pttm/models"
)

var DB *gorm.DB

func ConnectDB() {
	database, err := gorm.Open(sqlite.Open("pttm.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.Task{}) // This line is crucial
	DB = database
}

func InitDatabase() {
	dbPath := "pttm.db" // or from env/config file

	// Check if file exists, optionally log it
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Creating new database...")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate your models
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Database connected & migrated.")
}
