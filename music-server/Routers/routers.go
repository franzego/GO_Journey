package routers

import (
	"net/http"

	handlers "github.com/franzego/music-server/Handlers"
	middleware "github.com/franzego/music-server/Middleware"
	"github.com/gorilla/mux"
)

func RouterFunc() *mux.Router {
	rou := mux.NewRouter()

	// Request -> Middleware -> Main Logic(signup/login)

	rou.Handle("/signup", middleware.Loggingmiddleware(http.HandlerFunc(handlers.SignupHandler))).Methods("POST")
	rou.Handle("/login", middleware.Loggingmiddleware(http.HandlerFunc(handlers.LoginHandler))).Methods("POST")

	return rou
}
