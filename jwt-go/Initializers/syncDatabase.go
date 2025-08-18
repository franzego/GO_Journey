package initializers

import (
	models "github.com/franzego/jwt-go/Models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})

}
