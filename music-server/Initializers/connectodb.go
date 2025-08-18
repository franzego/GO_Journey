package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// This is the function that connects the gorm to the aiven-postgres(alternative to elephant sql) that i am using for this project.
// It ultimately creates the DB based on the gorm model that has been specified in models.go

var DB *gorm.DB

func Connectiontodb() {
	var err error

	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to create database")
	}
}
