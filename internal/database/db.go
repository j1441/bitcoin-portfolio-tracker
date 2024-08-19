package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	var db *sql.DB
	var err error

	// Get the database URL from the environment variable
	databaseUrl := os.Getenv("postgres://u8jks8baoetphi:p0c2176fb806d0f9b73959a75f12006ae4829cf910bfef8b5b04f8138dd177123@ceqbglof0h8enj.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com:5432/dat3dvea3hjegd")

	if databaseUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to the database
	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
