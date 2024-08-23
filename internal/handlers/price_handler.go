package handlers

import (
	"bitcoin-portfolio-tracker/internal/coingecko"
	"encoding/json"
	"log"
	"net/http"
)

func GetBitcoinPriceHandler(w http.ResponseWriter, r *http.Request) {
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
