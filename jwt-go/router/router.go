package router

import (
	handlers "github.com/franzego/jwt-go/Handlers"
	"github.com/gin-gonic/gin"
)

func RouterFunc() {
	r := gin.Default()

	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.Run() // listen and serve on 0.0.0.0:8080
}
