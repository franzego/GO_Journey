package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lpernett/godotenv"
)

// This is the middleware handling the Authentication. i.e JWT token to ensure authentication

func Authmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, falling back to system environment")
		}
		token := r.Header.Get("Authorization")

		// Always strip "Bearer " prefix if present
		token = strings.TrimPrefix(token, "Bearer ")

		if token != os.Getenv("SECRET") {
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})

	//validator.New
}
