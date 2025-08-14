package server

import (
	"log"
	"net/http"
	"time"
)

// Logging - логирует запросы к серверу
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method: %s, URI: %s, time: %s", r.Method, r.RequestURI, time.Since(start))
	})
}
