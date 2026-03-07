# Repository Guidelines

## Project Structure & Module Organization
- `cmd/main.go`: CLI entrypoint and flag parsing (`-search`, `-token`, `-limit`, `-freq`, `-send`, `-help`).
- `exec/`: application actions (search, Telegram notifications, sell execution, transaction sending).
- `http/requests.go`: CoinGecko price-fetch layer.
- `config/config.go`: environment bootstrap (`.env` loading).
- `logs/`: reserved for logging utilities (currently minimal).
- Root files: `go.mod`, `go.sum`, `.env` (local secrets only), `Readme.md`.

## Build, Test, and Development Commands
- `go run ./cmd -help`: run the CLI and print available flags.
- `go run ./cmd -search avax`: search CoinGecko token IDs by name/symbol.
- `go run ./cmd -token ethereum -limit 3000 -freq 1h`: monitor price and trigger sell logic.
- `go run ./cmd -send`: send a Sepolia test transaction using `.env` credentials.
- `go build ./cmd`: compile the application.
- `go test ./...`: run all tests (currently expected to be sparse until tests are added).
- `go fmt ./... && go vet ./...`: format and run static checks before opening a PR.

## Coding Style & Naming Conventions
- Follow standard Go formatting (`gofmt`) and idioms.
- Package names are short, lowercase nouns (`exec`, `config`, `http`).
- Exported identifiers use `PascalCase`; internal helpers use `camelCase`.
- Keep functions focused; isolate API calls and side effects from CLI flag handling.
- Prefer English for new code/comments to keep a single language style across contributors.

## Testing Guidelines
- Place tests next to code as `*_test.go` files (table-driven tests preferred).
- Mock external calls (CoinGecko, Telegram, RPC) instead of hitting live endpoints in unit tests.
- Cover parsing, error paths, and boundary values (e.g., invalid `-freq`, missing env vars).
- Run `go test ./...` locally before pushing.

## Commit & Pull Request Guidelines
- Use short, imperative commit messages (e.g., `Add token price retry on HTTP 429`).
- Keep commits focused (one logical change per commit).
- PRs should include: purpose, key changes, how to run/verify, and related issue links.
- For CLI behavior changes, include sample command/output snippets.
- Never include real private keys, bot tokens, or chat IDs in commits or screenshots.

## Security & Configuration Tips
- Keep secrets in `.env`; do not commit `.env` or raw credentials.
- Use testnet values for transaction development (`RPC_URL`, `PRIVATE_KEY`, `DESTINATION_ACCOUNT`).
- Validate required env vars at startup for features that depend on external services.
