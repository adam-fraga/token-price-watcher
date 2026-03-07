package exec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type TelegramUpdateResponse struct {
	OK     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

type TelegramUpdate struct {
	UpdateID    int64           `json:"update_id"`
	Message     *TelegramRecord `json:"message,omitempty"`
	ChannelPost *TelegramRecord `json:"channel_post,omitempty"`
}

type TelegramRecord struct {
	MessageID int64  `json:"message_id"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
	Chat      struct {
		ID int64 `json:"id"`
	} `json:"chat"`
}

type ParsedSignal struct {
	Action string
	Token  string
	Entry  *float64
	Target *float64
	Stop   *float64
	Raw    string
}

func FetchTelegramUpdates(offset int64, limit int) ([]TelegramUpdate, error) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is empty")
	}

	if limit <= 0 {
		limit = 50
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=1&offset=%d&limit=%d", botToken, offset, limit)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed TelegramUpdateResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	if !parsed.OK {
		return nil, fmt.Errorf("telegram getUpdates failed")
	}
	return parsed.Result, nil
}

func LoadTelegramOffset(path string) (int64, error) {
	if path == "" {
		return 0, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	raw := strings.TrimSpace(string(b))
	if raw == "" {
		return 0, nil
	}
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func SaveTelegramOffset(path string, offset int64) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(fmt.Sprintf("%d\n", offset)), 0o644)
}

func ParseSignalText(text string) (ParsedSignal, bool) {
	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return ParsedSignal{}, false
	}
	upper := strings.ToUpper(normalized)

	action := ""
	switch {
	case strings.Contains(upper, "BUY"), strings.Contains(upper, "LONG"):
		action = "BUY"
	case strings.Contains(upper, "SELL"), strings.Contains(upper, "SHORT"):
		action = "SELL"
	case strings.Contains(upper, "CLOSE"), strings.Contains(upper, "TAKE PROFIT"):
		action = "CLOSE"
	case strings.Contains(upper, "CANCEL"), strings.Contains(upper, "INVALID"):
		action = "CANCEL"
	}
	if action == "" {
		return ParsedSignal{}, false
	}

	token := extractToken(upper)
	if token == "" {
		return ParsedSignal{}, false
	}

	signal := ParsedSignal{
		Action: action,
		Token:  token,
		Entry:  extractPrice(upper, `(?:ENTRY|ENTER|AT)\s*[:=]?\s*([0-9]+(?:\.[0-9]+)?)`),
		Target: extractPrice(upper, `(?:TARGET|TP)\s*[:=]?\s*([0-9]+(?:\.[0-9]+)?)`),
		Stop:   extractPrice(upper, `(?:STOP|SL)\s*[:=]?\s*([0-9]+(?:\.[0-9]+)?)`),
		Raw:    normalized,
	}
	return signal, true
}

func extractToken(upper string) string {
	quoteMatcher := regexp.MustCompile(`\b([A-Z]{2,12})(USDT|USD|USDC|BTC|ETH|AVAX|SOL)\b`)
	if m := quoteMatcher.FindStringSubmatch(upper); len(m) > 0 {
		return m[0]
	}

	tickerMatcher := regexp.MustCompile(`\b([A-Z]{2,10})\b`)
	ignore := map[string]bool{
		"BUY": true, "LONG": true, "SELL": true, "SHORT": true, "CLOSE": true, "CANCEL": true,
		"ENTRY": true, "TARGET": true, "STOP": true, "SL": true, "TP": true, "AT": true,
	}
	matches := tickerMatcher.FindAllString(upper, -1)
	for _, m := range matches {
		if !ignore[m] {
			return m
		}
	}
	return ""
}

func extractPrice(text string, pattern string) *float64 {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return nil
	}
	val, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return nil
	}
	return &val
}

func FormatSignalNotification(sourceChatID string, signal ParsedSignal) string {
	var details []string
	details = append(details, fmt.Sprintf("action=%s", signal.Action))
	details = append(details, fmt.Sprintf("token=%s", signal.Token))
	if signal.Entry != nil {
		details = append(details, fmt.Sprintf("entry=%.4f", *signal.Entry))
	}
	if signal.Target != nil {
		details = append(details, fmt.Sprintf("target=%.4f", *signal.Target))
	}
	if signal.Stop != nil {
		details = append(details, fmt.Sprintf("stop=%.4f", *signal.Stop))
	}

	raw := signal.Raw
	if len(raw) > 140 {
		raw = raw[:140] + "..."
	}
	return fmt.Sprintf("📥 Parsed signal from chat %s | %s | raw=%s", sourceChatID, strings.Join(details, " "), raw)
}
