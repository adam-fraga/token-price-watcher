package exec

import "testing"

func TestParseSignalText(t *testing.T) {
	s, ok := ParseSignalText("BUY AVAXUSDT entry 21.5 target 25 stop 19.8")
	if !ok {
		t.Fatalf("expected parsed signal")
	}
	if s.Action != "BUY" || s.Token != "AVAXUSDT" {
		t.Fatalf("unexpected parsed values: %#v", s)
	}
	if s.Entry == nil || *s.Entry != 21.5 {
		t.Fatalf("expected entry=21.5, got %#v", s.Entry)
	}
	if s.Target == nil || *s.Target != 25 {
		t.Fatalf("expected target=25, got %#v", s.Target)
	}
	if s.Stop == nil || *s.Stop != 19.8 {
		t.Fatalf("expected stop=19.8, got %#v", s.Stop)
	}
}

func TestParseSignalTextNoAction(t *testing.T) {
	if _, ok := ParseSignalText("watching market conditions only"); ok {
		t.Fatalf("expected no signal")
	}
}
