package main

import (
	"bitcoin-portfolio-tracker/internal/database"
	"bitcoin-portfolio-tracker/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting app")

	// Initialize the database connection
	db := database.ConnectDB()
	defer db.Close()

	log.Println("Connected to db")

	http.HandleFunc("/signup", handlers.SignUpHandler(db))
	http.HandleFunc("/login", handlers.LoginHandler(db))
	http.HandleFunc("/portfolio", handlers.CreatePortfolioHandler(db))
	http.HandleFunc("/portfolios", handlers.ListPortfoliosHandler(db))
	http.HandleFunc("/portfolio/delete", handlers.DeletePortfolioHandler(db))

	http.HandleFunc("/price", handlers.GetBitcoinPriceHandler)

	// New route for M2 change data
	http.HandleFunc("/api/m2-change", handlers.GetM2ChangeHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
