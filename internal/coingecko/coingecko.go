package coingecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const apiURL = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"

type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
	} `json:"bitcoin"`
}

// FetchBitcoinPrice fetches the current Bitcoin price in USD from CoinGecko
func FetchBitcoinPrice() (float64, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch bitcoin price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("unexpected status code from CoinGecko API")
	}

	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode CoinGecko response: %w", err)
	}

	return result.Bitcoin.USD, nil
}
