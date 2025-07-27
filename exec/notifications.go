package exec

import (
	"fmt"
	cfg "github.com/adam-fraga/token-price-watcher/config"
	"net/http"
	"net/url"
	"os"
)

func SendTPWBotNotification(msg string) error {

	cfg.InitData()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	escapedMsg := url.QueryEscape(msg)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botToken, chatID, escapedMsg)

	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("❌ Erreur Telegram API, code: %d", resp.StatusCode)
	}
	return nil
}
