package main

import (
	"bitcoin_portfolio_tracker/internal/database"
	"bitcoin_portfolio_tracker/internal/handlers"
	"log"
	"net/http"
	//"github.com/rs/cors"
)

func main() {
	// Initialize the database connection
	db := database.ConnectDB()
	defer db.Close()

	//mux := http.NewServeMux()

	http.HandleFunc("/signup", handlers.SignUpHandler(db))
	http.HandleFunc("/login", handlers.LoginHandler(db))
	http.HandleFunc("/portfolio", handlers.CreatePortfolioHandler(db))
	http.HandleFunc("/portfolios", handlers.ListPortfoliosHandler(db))
	http.HandleFunc("/portfolio/delete", handlers.DeletePortfolioHandler(db))

	http.HandleFunc("/price", handlers.GetBitcoinPriceHandler)

	//handler := cors.Default().Handler(mux)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
