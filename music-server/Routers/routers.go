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
	rou.Handle("/oauth", middleware.Loggingmiddleware(http.HandlerFunc(handlers.OauthLoginHandler))).Methods("GET")
	//

	//Using subrouters in mux. This is so that there are public and protected routes
	api := rou.PathPrefix("/api").Subrouter()
	api.Use(middleware.Authmiddleware)

	//api.HandleFunc("/upload", handlers.UploadHandler).Methods("POST") // :8080/api/upload
	api.HandleFunc("/upload", handlers.TrackUploader).Methods("POST") // :8080/api/upload

	return rou
}
