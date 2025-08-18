package middleware

import (
	"log"
	"net/http"
	"time"
)

// This is the middleware handling the Logging

func Loggingmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("%s %s took %v", r.Method, r.URL.Path, duration)
	})
}
