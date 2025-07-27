package main

import (
	"flag"
	"fmt"
	"time"

	ex "github.com/adam-fraga/token-price-watcher/exec"
	h "github.com/adam-fraga/token-price-watcher/http"
)

func main() {
	help := flag.Bool("help", false, "📘 Afficher l'aide")
	searchQuery := flag.String("search", "", "Rechercher l'ID CoinGecko d'un token par nom ou symbole")
	tokenID := flag.String("token", "ethereum", "ID CoinGecko du token")
	limit := flag.Float64("limit", 3000.0, "Prix limite pour vendre")
	freq := flag.String("freq", "1h", "Fréquence (ex: 10s, 5m, 1h)")
	flag.Parse()

	if *help {
		ex.PrintHelp()
		return
	}

	if *searchQuery != "" {
		ex.SearchTokenID(*searchQuery)
		return
	}

	// Charger la config seulement si pas en mode search
	// cfg := config.Load()

	// Convertir la fréquence en durée
	duration, err := time.ParseDuration(*freq)
	if err != nil {
		fmt.Printf("Erreur dans la fréquence : %v\n", err)
		return
	}

	ex.SendTPWBotNotification(fmt.Sprintf("⏱️ Surveillance de %s chaque %s, vente si > %.2f USD", *tokenID, duration, *limit))

	for {
		price, err := h.GetTokenPrice(*tokenID)
		if err != nil {
			fmt.Printf("❌ Erreur récupération prix : %v\n", err)
		} else {
			ex.SendTPWBotNotification(fmt.Sprintf("📊 Prix de %s : %.2f USD", *tokenID, price))
			if price >= *limit {
				ex.ExecuteSell(*tokenID, price)
				break
			}
		}
		time.Sleep(duration)
	}
}
