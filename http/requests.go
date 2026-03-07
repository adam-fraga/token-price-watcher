package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"os"
	"strings"
	"time"
)

type CoinGeckoProvider struct {
	Client *stdhttp.Client
}

func (p *CoinGeckoProvider) Name() string {
	return "coingecko"
}

func (p *CoinGeckoProvider) GetTokenPrice(tokenID string) (float64, error) {
	client := p.Client
	if client == nil {
		client = &stdhttp.Client{Timeout: 10 * time.Second}
	}

	baseURL := os.Getenv("COINGECKO_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.coingecko.com/api/v3"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", baseURL, tokenID)
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != stdhttp.StatusOK {
		return 0, fmt.Errorf("coingecko returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	tokenPrice, ok := result[tokenID]
	if !ok {
		return 0, errors.New("token not found in provider response")
	}
	price, ok := tokenPrice["usd"]
	if !ok {
		return 0, errors.New("usd price missing in provider response")
	}

	return price, nil
}

// GetTokenPrice is kept for backward compatibility and uses CoinGecko.
func GetTokenPrice(tokenID string) (float64, error) {
	return (&CoinGeckoProvider{}).GetTokenPrice(tokenID)
}
