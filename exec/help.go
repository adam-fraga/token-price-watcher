package exec

import "fmt"

func PrintHelp() {
	fmt.Print(`🛠️  Utilisation :

  -search <mot-clé>      → Recherche les tokens liés au nom ou symbole donné
  -token <id>            → Token (CoinGecko: avalanche-2, Binance: AVAXUSDT)
  -price-provider <nom>  → Source de prix (coingecko | binance)
  -freq <durée>          → Fréquence de polling (défaut: 12h)
  -notify-every <durée>  → Notification de statut périodique (défaut: 24h)
  -cooldown <durée>      → Cooldown entre alertes de seuil (défaut: 6h)
  -upper <prix>          → Alerte si prix >= upper (0 désactive)
  -lower <prix>          → Alerte si prix <= lower (0 désactive)
  -signal-ingest         → Active l'ingestion de signaux Telegram (getUpdates)
  -signal-chat-ids       → Chats/channels source CSV (sinon TELEGRAM_SOURCE_CHAT_IDS)
  -signal-poll-every     → Polling signaux Telegram (défaut: 30m)
  -signal-offset-file    → Fichier offset signaux (sinon env/default)
  -once                  → Envoie une seule notification puis quitte
  -help                 → Affiche cette aide

📌 Exemples :
  go run ./cmd -search avax
  go run ./cmd -token avalanche-2 -price-provider coingecko -freq 12h -notify-every 24h -cooldown 6h -upper 18 -lower 12
  go run ./cmd -token AVAXUSDT -price-provider binance -once
  go run ./cmd -token avalanche-2 -signal-ingest -signal-chat-ids "-100123,-100456" -signal-poll-every 15m
`)
}
