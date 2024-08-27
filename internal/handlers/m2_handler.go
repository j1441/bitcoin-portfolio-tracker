package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func GetM2ChangeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		addCORSHeaders(w, r)

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log the request
		log.Println("Received request on /m2-change")

		// Check for recent M2 change data in the database
		var m2ChangeValue float64
		var timestamp time.Time
		err := db.QueryRow("SELECT m2_change_value, timestamp FROM m2_data ORDER BY timestamp DESC LIMIT 1").Scan(&m2ChangeValue, &timestamp)

		if err != nil {
			if err != sql.ErrNoRows {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			// Fetch new data if no recent data exists
			m2ChangeValue, err = fetchM2ChangeFromAPI(db)
			if err != nil {
				http.Error(w, "Failed to fetch M2 data", http.StatusInternalServerError)
				return
			}
		} else if time.Since(timestamp).Hours() >= 24 {
			// Fetch new data if the cached data is older than 24 hours
			m2ChangeValue, err = fetchM2ChangeFromAPI(db)
			if err != nil {
				http.Error(w, "Failed to fetch M2 data", http.StatusInternalServerError)
				return
			}
		}

		// Return the M2 change value
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"m2_change": m2ChangeValue})
	}
}

// fetchM2ChangeFromAPI fetches the M2 change value from the SSB API and stores it in the database
func fetchM2ChangeFromAPI(db *sql.DB) (float64, error) {
	resp, err := http.Get("https://data.ssb.no/api/v0/dataset/172769.json?lang=no")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var parsedResponse struct {
		Dataset struct {
			Value []float64 `json:"value"`
		} `json:"dataset"`
	}

	err = json.NewDecoder(resp.Body).Decode(&parsedResponse)
	if err != nil {
		return 0, err
	}

	// Assume the last value is the M2 change
	m2ChangeValue := parsedResponse.Dataset.Value[len(parsedResponse.Dataset.Value)-1]

	// Store the new M2 change value in the database
	_, err = db.Exec("INSERT INTO m2_data (m2_change_value, timestamp) VALUES ($1, $2)", m2ChangeValue, time.Now())
	if err != nil {
		return 0, err
	}

	return m2ChangeValue, nil
}
