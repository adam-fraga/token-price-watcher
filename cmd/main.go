package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	cfg "github.com/adam-fraga/token-price-watcher/config"
	ex "github.com/adam-fraga/token-price-watcher/exec"
	h "github.com/adam-fraga/token-price-watcher/http"
)

func main() {
	help := flag.Bool("help", false, "📘 Afficher l'aide")
	searchQuery := flag.String("search", "", "Rechercher l'ID CoinGecko d'un token par nom ou symbole")
	tokenID := flag.String("token", "avalanche-2", "Référence token (CoinGecko ID ou Binance symbol)")
	freq := flag.String("freq", "12h", "Fréquence (ex: 30s, 5m, 1h, 12h)")
	priceProviderName := flag.String("price-provider", "coingecko", "Source de prix (coingecko|binance)")
	upper := flag.Float64("upper", 0, "Alerte si prix >= upper (0 désactive)")
	lower := flag.Float64("lower", 0, "Alerte si prix <= lower (0 désactive)")
	notifyEvery := flag.String("notify-every", "24h", "Fréquence de notification de statut (ex: 6h, 24h)")
	cooldown := flag.String("cooldown", "6h", "Cooldown entre alertes de seuil (ex: 30m, 6h)")
	signalIngest := flag.Bool("signal-ingest", false, "Ingestion de signaux depuis Telegram getUpdates")
	signalChatIDs := flag.String("signal-chat-ids", "", "IDs de chats/channels source séparés par virgule (sinon TELEGRAM_SOURCE_CHAT_IDS)")
	signalPollEvery := flag.String("signal-poll-every", "30m", "Fréquence de polling des signaux Telegram")
	signalOffsetFile := flag.String("signal-offset-file", "", "Fichier offset Telegram (sinon TELEGRAM_SIGNAL_OFFSET_FILE ou data/telegram_signal_offset.txt)")
	once := flag.Bool("once", false, "Envoie une notification unique puis quitte")
	flag.Parse()

	if *help {
		ex.PrintHelp()
		return
	}

	if *searchQuery != "" {
		ex.SearchTokenID(*searchQuery)
		return
	}

	// Charge .env pour providers et Telegram.
	cfg.InitData()

	pollInterval, err := time.ParseDuration(*freq)
	if err != nil {
		fmt.Printf("Erreur dans la fréquence : %v\n", err)
		return
	}

	statusInterval, err := time.ParseDuration(*notifyEvery)
	if err != nil {
		fmt.Printf("Erreur notify-every : %v\n", err)
		return
	}
	cooldownInterval, err := time.ParseDuration(*cooldown)
	if err != nil {
		fmt.Printf("Erreur cooldown : %v\n", err)
		return
	}
	signalPollInterval, err := time.ParseDuration(*signalPollEvery)
	if err != nil {
		fmt.Printf("Erreur signal-poll-every : %v\n", err)
		return
	}

	priceProvider, err := h.NewPriceProvider(*priceProviderName)
	if err != nil {
		fmt.Printf("❌ Erreur provider prix : %v\n", err)
		return
	}

	startMsg := fmt.Sprintf("🔔 TPW notifications started | token=%s | provider=%s | poll=%s | status=%s | upper=%.4f | lower=%.4f",
		*tokenID, priceProvider.Name(), pollInterval, statusInterval, *upper, *lower)
	fmt.Println(startMsg)
	_ = ex.SendTPWBotNotification(startMsg)

	if !*signalIngest && strings.EqualFold(os.Getenv("TELEGRAM_SIGNAL_ENABLED"), "true") {
		*signalIngest = true
	}

	if *signalOffsetFile == "" {
		*signalOffsetFile = os.Getenv("TELEGRAM_SIGNAL_OFFSET_FILE")
		if *signalOffsetFile == "" {
			*signalOffsetFile = "data/telegram_signal_offset.txt"
		}
	}

	sourceChatIDsRaw := strings.TrimSpace(*signalChatIDs)
	if sourceChatIDsRaw == "" {
		sourceChatIDsRaw = strings.TrimSpace(os.Getenv("TELEGRAM_SOURCE_CHAT_IDS"))
	}

	sourceChatIDsMap := map[int64]bool{}
	if sourceChatIDsRaw != "" {
		parts := strings.Split(sourceChatIDsRaw, ",")
		for _, part := range parts {
			idStr := strings.TrimSpace(part)
			if idStr == "" {
				continue
			}
			id, parseErr := strconv.ParseInt(idStr, 10, 64)
			if parseErr != nil {
				fmt.Printf("❌ Chat ID invalide: %s (%v)\n", idStr, parseErr)
				return
			}
			sourceChatIDsMap[id] = true
		}
	}

	var signalOffset int64
	if *signalIngest {
		offset, loadErr := ex.LoadTelegramOffset(*signalOffsetFile)
		if loadErr != nil {
			fmt.Printf("❌ Impossible de charger offset signal: %v\n", loadErr)
			return
		}
		signalOffset = offset
		_ = ex.SendTPWBotNotification(fmt.Sprintf("📡 Signal ingestion ON | chats=%s | poll=%s | offset=%d", sourceChatIDsRaw, signalPollInterval, signalOffset))
	}

	var wasAboveUpper bool
	var wasBelowLower bool
	var lastUpperAlertAt time.Time
	var lastLowerAlertAt time.Time
	nextStatusAt := time.Now().Add(statusInterval)
	nextSignalPollAt := time.Now()

	for {
		price, err := priceProvider.GetTokenPrice(*tokenID)
		if err != nil {
			fmt.Printf("❌ Erreur récupération prix : %v\n", err)
			_ = ex.SendTPWBotNotification(fmt.Sprintf("❌ Price fetch failed for %s (%s): %v", *tokenID, priceProvider.Name(), err))
		} else {
			fmt.Printf("📊 %s: %.4f USD\n", *tokenID, price)
			now := time.Now()

			if *upper > 0 {
				allowedByCooldown := cooldownInterval <= 0 || lastUpperAlertAt.IsZero() || now.Sub(lastUpperAlertAt) >= cooldownInterval
				if price >= *upper && !wasAboveUpper && allowedByCooldown {
					_ = ex.SendTPWBotNotification(fmt.Sprintf("🚨 %s crossed ABOVE %.4f (current: %.4f)", *tokenID, *upper, price))
					wasAboveUpper = true
					lastUpperAlertAt = now
				}
				if price < *upper {
					wasAboveUpper = false
				}
			}

			if *lower > 0 {
				allowedByCooldown := cooldownInterval <= 0 || lastLowerAlertAt.IsZero() || now.Sub(lastLowerAlertAt) >= cooldownInterval
				if price <= *lower && !wasBelowLower && allowedByCooldown {
					_ = ex.SendTPWBotNotification(fmt.Sprintf("🚨 %s crossed BELOW %.4f (current: %.4f)", *tokenID, *lower, price))
					wasBelowLower = true
					lastLowerAlertAt = now
				}
				if price > *lower {
					wasBelowLower = false
				}
			}

			if now.After(nextStatusAt) || *once {
				_ = ex.SendTPWBotNotification(fmt.Sprintf("📊 Status %s (%s): %.4f USD", *tokenID, priceProvider.Name(), price))
				nextStatusAt = now.Add(statusInterval)
			}

			if *once {
				fmt.Println("✅ Notification sent once, exiting.")
				return
			}
		}

		now := time.Now()
		if *signalIngest && (now.After(nextSignalPollAt) || now.Equal(nextSignalPollAt)) {
			updates, fetchErr := ex.FetchTelegramUpdates(signalOffset+1, 100)
			if fetchErr != nil {
				fmt.Printf("❌ Signal ingestion error: %v\n", fetchErr)
				_ = ex.SendTPWBotNotification(fmt.Sprintf("❌ Signal ingestion error: %v", fetchErr))
			} else {
				for _, update := range updates {
					if update.UpdateID > signalOffset {
						signalOffset = update.UpdateID
					}

					rec := update.ChannelPost
					if rec == nil {
						rec = update.Message
					}
					if rec == nil || strings.TrimSpace(rec.Text) == "" {
						continue
					}

					if len(sourceChatIDsMap) > 0 && !sourceChatIDsMap[rec.Chat.ID] {
						continue
					}

					signal, ok := ex.ParseSignalText(rec.Text)
					if !ok {
						continue
					}

					msg := ex.FormatSignalNotification(strconv.FormatInt(rec.Chat.ID, 10), signal)
					fmt.Println(msg)
					_ = ex.SendTPWBotNotification(msg)
				}
				_ = ex.SaveTelegramOffset(*signalOffsetFile, signalOffset)
			}
			nextSignalPollAt = now.Add(signalPollInterval)
		}

		time.Sleep(pollInterval)
	}
}
