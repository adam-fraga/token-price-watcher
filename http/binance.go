package http

import (
	"encoding/json"
	"fmt"
	"io"
	stdhttp "net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type BinanceProvider struct {
	Client *stdhttp.Client
}

func (p *BinanceProvider) Name() string {
	return "binance"
}

func (p *BinanceProvider) GetTokenPrice(tokenRef string) (float64, error) {
	symbol := strings.ToUpper(strings.TrimSpace(tokenRef))
	if symbol == "" {
		return 0, fmt.Errorf("binance provider requires a market symbol, e.g. AVAXUSDT")
	}

	client := p.Client
	if client == nil {
		client = &stdhttp.Client{Timeout: 10 * time.Second}
	}

	baseURL := os.Getenv("BINANCE_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.binance.com"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", baseURL, symbol)
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != stdhttp.StatusOK {
		return 0, fmt.Errorf("binance returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid binance price %q: %w", result.Price, err)
	}
	return price, nil
}
