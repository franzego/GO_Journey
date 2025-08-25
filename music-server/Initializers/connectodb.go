package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// This is the function that connects the gorm to the aiven-postgres(alternative to elephant sql) that i am using for this project.
// It ultimately creates the DB based on the gorm model that has been specified in models.go

var DB *gorm.DB

func Connectiontodb() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment")
		return
	}
	//var err error
	// Load .env file
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Attempting to connect with DSN:", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	//DB = db
}
