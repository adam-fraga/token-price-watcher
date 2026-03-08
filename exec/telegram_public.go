package exec

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type TelegramPublicPost struct {
	Channel string
	PostID  int64
	Text    string
	URL     string
}

var (
	telegramPublicPostRefRE = regexp.MustCompile(`data-post="([^"/]+)/([0-9]+)"`)
	telegramMessageTextRE   = regexp.MustCompile(`(?s)<div class="tgme_widget_message_text[^"]*"[^>]*>(.*?)</div>`)
	telegramBRTagRE         = regexp.MustCompile(`(?i)<br\s*/?>`)
	telegramHTMLTagRE       = regexp.MustCompile(`(?s)<[^>]+>`)
)

func FetchTelegramPublicPosts(channelRef string) ([]TelegramPublicPost, error) {
	channel, pageURL, err := NormalizeTelegramPublicChannel(channelRef)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram public page returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return ParseTelegramPublicPosts(string(body), channel), nil
}

func NormalizeTelegramPublicChannel(channelRef string) (string, string, error) {
	normalized := strings.TrimSpace(channelRef)
	normalized = strings.TrimSuffix(normalized, "/")
	normalized = strings.TrimPrefix(normalized, "@")
	normalized = strings.TrimPrefix(normalized, "https://t.me/s/")
	normalized = strings.TrimPrefix(normalized, "http://t.me/s/")
	normalized = strings.TrimPrefix(normalized, "https://t.me/")
	normalized = strings.TrimPrefix(normalized, "http://t.me/")
	normalized = strings.TrimPrefix(normalized, "https://telegram.me/s/")
	normalized = strings.TrimPrefix(normalized, "http://telegram.me/s/")
	normalized = strings.TrimPrefix(normalized, "https://telegram.me/")
	normalized = strings.TrimPrefix(normalized, "http://telegram.me/")
	normalized = strings.Trim(normalized, "/")

	if normalized == "" {
		return "", "", fmt.Errorf("telegram public channel is empty")
	}

	return normalized, "https://t.me/s/" + normalized, nil
}

func ParseTelegramPublicPosts(pageHTML string, fallbackChannel string) []TelegramPublicPost {
	matches := telegramPublicPostRefRE.FindAllStringSubmatchIndex(pageHTML, -1)
	if len(matches) == 0 {
		return nil
	}

	posts := make([]TelegramPublicPost, 0, len(matches))
	for i, match := range matches {
		if len(match) < 6 {
			continue
		}

		start := match[0]
		end := len(pageHTML)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		block := pageHTML[start:end]

		channel := pageHTML[match[2]:match[3]]
		if channel == "" {
			channel = fallbackChannel
		}

		postID, err := strconv.ParseInt(pageHTML[match[4]:match[5]], 10, 64)
		if err != nil {
			continue
		}

		text := ""
		if textMatch := telegramMessageTextRE.FindStringSubmatch(block); len(textMatch) > 1 {
			text = cleanTelegramHTMLText(textMatch[1])
		}

		posts = append(posts, TelegramPublicPost{
			Channel: channel,
			PostID:  postID,
			Text:    text,
			URL:     fmt.Sprintf("https://t.me/%s/%d", channel, postID),
		})
	}

	return posts
}

func LatestParsedSignalFromPosts(posts []TelegramPublicPost) (TelegramPublicPost, ParsedSignal, bool) {
	var latestPost TelegramPublicPost
	var latestSignal ParsedSignal
	found := false

	for _, post := range posts {
		if strings.TrimSpace(post.Text) == "" {
			continue
		}

		signal, ok := ParseSignalText(post.Text)
		if !ok {
			continue
		}

		if !found || post.PostID > latestPost.PostID {
			latestPost = post
			latestSignal = signal
			found = true
		}
	}

	return latestPost, latestSignal, found
}

func LatestParsedPublicSignal(posts []TelegramPublicPost) (TelegramPublicPost, ParsedSignal, bool) {
	return LatestParsedSignalFromPosts(posts)
}

func cleanTelegramHTMLText(fragment string) string {
	fragment = telegramBRTagRE.ReplaceAllString(fragment, "\n")
	fragment = telegramHTMLTagRE.ReplaceAllString(fragment, " ")
	fragment = html.UnescapeString(fragment)

	lines := strings.Split(fragment, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.Join(strings.Fields(line), " ")
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}

	return strings.Join(cleaned, "\n")
}
