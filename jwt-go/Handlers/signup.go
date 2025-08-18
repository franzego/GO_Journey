package handlers

import (
	"net/http"

	initializers "github.com/franzego/jwt-go/Initializers"
	models "github.com/franzego/jwt-go/Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//Get email/password

	type Required struct {
		Email    string
		Password string
	}
	var body Required

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(201, body)

	//hash password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Failed to hash password": err.Error()})
		return
	}

	//create user

	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{})
}
