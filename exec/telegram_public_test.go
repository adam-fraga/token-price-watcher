package exec

import "testing"

func TestNormalizeTelegramPublicChannel(t *testing.T) {
	channel, pageURL, err := NormalizeTelegramPublicChannel("https://t.me/s/tofan_trade/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if channel != "tofan_trade" {
		t.Fatalf("expected channel tofan_trade, got %q", channel)
	}

	if pageURL != "https://t.me/s/tofan_trade" {
		t.Fatalf("unexpected page URL: %q", pageURL)
	}
}

func TestParseTelegramPublicPostsAndLatestSignal(t *testing.T) {
	pageHTML := `
<div class="tgme_widget_message_wrap">
  <div class="tgme_widget_message" data-post="tofan_trade/41">
    <div class="tgme_widget_message_text js-message_text">Market update only</div>
  </div>
</div>
<div class="tgme_widget_message_wrap">
  <div class="tgme_widget_message" data-post="tofan_trade/42">
    <div class="tgme_widget_message_text js-message_text">BUY AVAXUSDT<br>Entry: 21.5<br>Target: 25<br>SL: 19.8</div>
  </div>
</div>
<div class="tgme_widget_message_wrap">
  <div class="tgme_widget_message" data-post="tofan_trade/43">
    <div class="tgme_widget_message_text js-message_text">SELL BTCUSDT<br>AT 63500<br>TP 62000<br>STOP 64200</div>
  </div>
</div>`

	posts := ParseTelegramPublicPosts(pageHTML, "tofan_trade")
	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	post, signal, ok := LatestParsedPublicSignal(posts)
	if !ok {
		t.Fatalf("expected latest parsed signal")
	}

	if post.PostID != 43 {
		t.Fatalf("expected latest post ID 43, got %d", post.PostID)
	}

	if signal.Action != "SELL" || signal.Token != "BTCUSDT" {
		t.Fatalf("unexpected signal: %#v", signal)
	}

	if signal.Entry == nil || *signal.Entry != 63500 {
		t.Fatalf("expected entry=63500, got %#v", signal.Entry)
	}
	if signal.Target == nil || *signal.Target != 62000 {
		t.Fatalf("expected target=62000, got %#v", signal.Target)
	}
	if signal.Stop == nil || *signal.Stop != 64200 {
		t.Fatalf("expected stop=64200, got %#v", signal.Stop)
	}
}
