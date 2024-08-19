package handlers

import (
	"bitcoin_portfolio_tracker/internal/coingecko"
	"encoding/json"
	"log"
	"net/http"
)

func GetBitcoinPriceHandler(w http.ResponseWriter, r *http.Request) {
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
