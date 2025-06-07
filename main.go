package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// Token attendu (à définir dans une variable d'environnement)
var expectedToken = os.Getenv("CONFIG_TOKEN")

func main() {
	if expectedToken == "" {
		log.Fatal("CONFIG_TOKEN is not set")
	}

	http.HandleFunc("/config", authMiddleware(configHandler))

	port := ":8080"
	log.Printf("Config server running on %s", port)
	log.Fatal(http.ListenAndServe(port, nil)) // pour HTTPS: ListenAndServeTLS
}

// Middleware d'authentification
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Unauthorized access attempt")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token != expectedToken {
			log.Println("Forbidden access attempt with token:", token)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Handler pour envoyer le fichier properties
func configHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "config/application.yaml"
	http.ServeFile(w, r, filePath)
}
