package handlers

import (
	"net/http"
	"os"
	"time"

	initializers "github.com/franzego/jwt-go/Initializers"
	models "github.com/franzego/jwt-go/Models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	//get email and password

	type Required struct {
		Email    string
		Password string
	}
	var body Required

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//look up reqested user
	var user models.User

	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	//compare the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	//generate jwt

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.ID,
			"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldnt generate valid jwt token"})
		return
	}

	//sWe set up cookie
	c.SetCookie(
		"jwt",       // cookie name
		tokenString, // value
		3600*24,     // max age in seconds (1 day)
		"/",         // path
		"localhost", // domain (adjust in prod)
		true,        // secure (true in prod with HTTPS)
		true,        // httpOnly
	)
}
