package http

import "testing"

func TestNewPriceProvider(t *testing.T) {
	p, err := NewPriceProvider("coingecko")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatalf("expected provider, got nil")
	}
	if p.Name() != "coingecko" {
		t.Fatalf("unexpected provider name: %s", p.Name())
	}
}

func TestNewPriceProviderBinance(t *testing.T) {
	p, err := NewPriceProvider("binance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatalf("expected provider, got nil")
	}
	if p.Name() != "binance" {
		t.Fatalf("unexpected provider name: %s", p.Name())
	}
}

func TestNewPriceProviderUnsupported(t *testing.T) {
	if _, err := NewPriceProvider("unknown"); err == nil {
		t.Fatalf("expected error for unsupported provider")
	}
}
