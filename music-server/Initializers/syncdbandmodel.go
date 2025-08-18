package initializers

import models "github.com/franzego/music-server/Models"

//This syncs the DB to the models

func Syncdbandmodels() {
	DB.AutoMigrate(&models.User{}, &models.Album{}, &models.Track{})
}
