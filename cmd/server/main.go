package main

import (
	"bitcoin-portfolio-tracker/internal/database"
	"bitcoin-portfolio-tracker/internal/handlers"
	"log"
	"net/http"

	//"github.com/rs/cors"
	//test

	"os"
)

func main() {
	// Initialize the database connection
	db := database.ConnectDB()
	defer db.Close()

	log.Println("Connected to db")

	//mux := http.NewServeMux()

	http.HandleFunc("/signup", handlers.SignUpHandler(db))
	http.HandleFunc("/login", handlers.LoginHandler(db))
	http.HandleFunc("/portfolio", handlers.CreatePortfolioHandler(db))
	http.HandleFunc("/portfolios", handlers.ListPortfoliosHandler(db))
	http.HandleFunc("/portfolio/delete", handlers.DeletePortfolioHandler(db))

	http.HandleFunc("/price", handlers.GetBitcoinPriceHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//handler := cors.Default().Handler(mux)

	log.Println("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
