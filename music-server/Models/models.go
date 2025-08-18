package models

import (
	"time"

	"gorm.io/gorm"
)

// simply creating models that will be used to manage data coming in

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `json:"-"`

	Tracks []Track
	Albums []Album
}

// struct for the response to a successful signup
type SignupResponse struct {
	Email string `json:"email"`
	Msg   string `json:"msg"`
}

// struct for album
type Album struct {
	ID        int `gorm:"primary key"`
	Title     string
	Artist    string
	UserID    uint `gorm:"not null"` //ownwer
	CreatedAt time.Time
	UpdatedAt time.Time
	Tracks    []Track
}

// struct for tracks
type Track struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	FilePath  string `gorm:"not null"` // path on server or cloud
	Duration  int    // duration in seconds
	AlbumID   uint   // optional, for tracks in an album
	UserID    uint   `gorm:"not null"` // uploader
	CreatedAt time.Time
	UpdatedAt time.Time
}
