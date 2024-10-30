package db

import (
	"Trecker/internal/db/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	dsn := "host=localhost user=zelmoron password=zelmoron443 dbname=habits port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Миграция схемы
	if err := db.AutoMigrate(&models.UsersModel{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	log.Println("Database migrated successfully")
	// Миграция схемы
	if err := db.AutoMigrate(&models.Habit{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")
}

func GetDB() *gorm.DB {
	return db
}
