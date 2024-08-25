// cors.go
package handlers

import (
	"net/http"
)

func isAllowedOrigin(origin string) bool {
	allowedOrigins := []string{
		"http://web.numerisgroup.xyz",       // Replace with your actual web frontend URL
		"http://another-allowed-origin.com", // Add other allowed origins if necessary
	}
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if origin == "" {
		// Allow mobile apps or other non-browser requests without an Origin header
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
