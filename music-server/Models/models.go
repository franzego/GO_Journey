package models

import "gorm.io/gorm"

// simply creating models that will be used to manage data coming in

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `json:"-"`
}

// struct for the response to a successful signup
type SignupResponse struct {
	Email string `json:"email"`
	Msg   string `json:"msg"`
}
