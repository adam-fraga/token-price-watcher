package http

import "fmt"

type PriceProvider interface {
	Name() string
	GetTokenPrice(tokenID string) (float64, error)
}

func NewPriceProvider(name string) (PriceProvider, error) {
	switch name {
	case "", "coingecko":
		return &CoinGeckoProvider{}, nil
	case "binance":
		return &BinanceProvider{}, nil
	default:
		return nil, fmt.Errorf("unsupported price provider: %s", name)
	}
}
