package handlers

import (
	"bitcoin-portfolio-tracker/internal/auth"
	"bitcoin-portfolio-tracker/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SignUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // or specify your frontend origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // If using cookies or authentication headers

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received request on /signup")
		var req models.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", req.Email, hashedPassword)
		if err != nil {
			log.Println("Error inserting user:", err) // Log the error to see what's wrong
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // or specify your frontend origin.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // If using cookies or authentication headers

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log the request to confirm it's reaching the handler
		log.Println("Received request on /login")

		var req models.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var userID int
		var hashedPassword string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE email=$1", req.Email).Scan(&userID, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if !auth.CheckPasswordHash(req.Password, hashedPassword) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Generate JWT
		tokenString, err := auth.GenerateJWT(userID, req.Email)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		fmt.Fprintf(w, "Login successful")
	}
}
