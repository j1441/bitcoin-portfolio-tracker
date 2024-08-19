package handlers

import (
	"bitcoin_portfolio_tracker/internal/auth"
	"bitcoin_portfolio_tracker/internal/coingecko"
	"bitcoin_portfolio_tracker/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func CreatePortfolioHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Portfolio

		// Parse the JWT token and extract the user ID
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Set the UserID from the token claims
		req.UserID = claims.UserID

		// Insert the portfolio into the database
		err = db.QueryRow("INSERT INTO portfolios (user_id, name, amount) VALUES ($1, $2, $3) RETURNING id",
			req.UserID, req.Name, req.Amount).Scan(&req.ID)
		if err != nil {
			log.Println("Error inserting portfolio:", err)
			http.Error(w, "Failed to create portfolio", http.StatusInternalServerError)
			return
		}

		// Return the created portfolio as a response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
	}
}

func ListPortfoliosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var portfolios []models.Portfolio

		// Parse the JWT token
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query("SELECT id, user_id, name, amount, created_at FROM portfolios WHERE user_id = $1", claims.UserID)
		if err != nil {
			http.Error(w, "Failed to retrieve portfolios", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Fetch the current Bitcoin price
		bitcoinPrice, err := coingecko.FetchBitcoinPrice()
		if err != nil {
			log.Println("Error fetching Bitcoin price:", err)
			http.Error(w, "Failed to fetch Bitcoin price", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var portfolio models.Portfolio
			if err := rows.Scan(&portfolio.ID, &portfolio.UserID, &portfolio.Name, &portfolio.Amount, &portfolio.CreatedAt); err != nil {
				http.Error(w, "Failed to scan portfolio", http.StatusInternalServerError)
				return
			}

			// Calculate the value in USD
			portfolio.ValueUSD = portfolio.Amount * bitcoinPrice
			portfolios = append(portfolios, portfolio)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error retrieving portfolios", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(portfolios)
	}
}

func DeletePortfolioHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID int `json:"id"`
		}

		// Parse the JWT token
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Delete the portfolio from the database
		result, err := db.Exec("DELETE FROM portfolios WHERE id = $1 AND user_id = $2", req.ID, claims.UserID)
		if err != nil {
			http.Error(w, "Failed to delete portfolio", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "Portfolio not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
