# Token Price Watcher (Notification Only)

TPW is a lightweight Go CLI that monitors token prices and sends Telegram notifications.

This version is intentionally notification-only:
- no auto-trading
- no execution providers
- no position memory engine

## Features
- Search token IDs on CoinGecko (`-search`).
- Monitor prices from:
  - `coingecko` (`-token avalanche-2`)
  - `binance` (`-token AVAXUSDT`)
- Threshold alerts:
  - `-upper` (alert when price crosses above)
  - `-lower` (alert when price crosses below)
- Alert cooldown (`-cooldown`) to avoid repeated threshold spam.
- Periodic status updates (`-notify-every`).
- One-shot notification mode (`-once`).
- Telegram signal ingestion:
  - read channel/chat posts via Bot API `getUpdates`
  - parse signal-like messages (`BUY/SELL/CLOSE/CANCEL`)
  - forward parsed signals as Telegram notifications

## Setup
1. Install dependencies:
```bash
go mod download
```
2. Create `.env`:
```env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
COINGECKO_API_BASE_URL=https://api.coingecko.com/api/v3
BINANCE_API_BASE_URL=https://api.binance.com
TELEGRAM_SIGNAL_ENABLED=true
TELEGRAM_SOURCE_CHAT_IDS=-1001234567890,-1002222222222
TELEGRAM_SIGNAL_OFFSET_FILE=data/telegram_signal_offset.txt
```

Notes:
- `TELEGRAM_SOURCE_CHAT_IDS` accepts comma-separated chat/channel IDs.
- Bot must be present in source channels and allowed to read messages.
- `TELEGRAM_SIGNAL_OFFSET_FILE` stores last processed update ID.

## Usage
Show help:
```bash
go run ./cmd -help
```

Search token IDs:
```bash
go run ./cmd -search avax
```

Monitor with alerts (CoinGecko):
```bash
go run ./cmd \
  -token avalanche-2 \
  -price-provider coingecko \
  -freq 12h \
  -notify-every 24h \
  -cooldown 6h \
  -upper 18 \
  -lower 12
```

Monitor with alerts (Binance):
```bash
go run ./cmd \
  -token AVAXUSDT \
  -price-provider binance \
  -freq 6h \
  -notify-every 12h \
  -upper 20 \
  -lower 10
```

Send one notification and exit:
```bash
go run ./cmd -token avalanche-2 -price-provider coingecko -once
```

Ingest signals from Telegram channels (using `.env` channels):
```bash
go run ./cmd \
  -token avalanche-2 \
  -price-provider coingecko \
  -signal-ingest \
  -signal-poll-every 15m
```

Override source channels from CLI:
```bash
go run ./cmd -signal-ingest -signal-chat-ids "-1001234567890,-1002222222222"
```

## Validation
```bash
gofmt -w ./...
go test ./...
```
