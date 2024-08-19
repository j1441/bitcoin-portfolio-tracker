package handlers

import (
	"bitcoin-portfolio-tracker/internal/coingecko"
	"encoding/json"
	"log"
	"net/http"
)

func GetBitcoinPriceHandler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Println("Received request on /price")

	price, err := coingecko.FetchBitcoinPrice()
	if err != nil {
		log.Println("Error fetching Bitcoin price:", err)
		http.Error(w, "Failed to fetch Bitcoin price", http.StatusInternalServerError)
		return
	}

	response := struct {
		Price float64 `json:"price"`
	}{
		Price: price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
