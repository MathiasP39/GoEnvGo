package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/joho/godotenv"
)



func main() {
	if os.Getenv("CONFIG_TOKEN") == "" {
		log.Printf("No CONFIG_TOKEN found in environment, loading from .env file");
		err := godotenv.Load()
		if err != nil {
    		log.Fatal("Error loading .env file")
  		}
	}
	var expectedToken = os.Getenv("CONFIG_TOKEN")
	if expectedToken == "" {
		log.Fatal("CONFIG_TOKEN is not set")
	}
	log.Printf("Using token: %s", expectedToken)
	log.Println("Starting config server...")
	http.HandleFunc("/config", authMiddleware(configHandler))
	http.HandleFunc("/config/{project_name}", authMiddleware(configHandler))

	port := ":8080"
	log.Printf("Config server running on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Unauthorized access attempt")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token != getExpectedToken() {
			log.Println("Forbidden access attempt with token:", token)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "config/application.yaml"
	http.ServeFile(w, r, filePath)
}

func getExpectedToken() string {
	expectedToken := os.Getenv("CONFIG_TOKEN")
	if expectedToken == "" {
		log.Fatal("CONFIG_TOKEN is not set")
	}
	return expectedToken
}
