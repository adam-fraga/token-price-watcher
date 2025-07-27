package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetTokenPrice(tokenID string) (float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", tokenID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	return result[tokenID]["usd"], nil
}
