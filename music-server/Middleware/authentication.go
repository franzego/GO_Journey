package middleware

import (
	"net/http"
	"os"
)

// This is the middleware handling the Authentication. i.e JWT token to ensure authentication

func Authmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token != os.Getenv("SECRET") {
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
