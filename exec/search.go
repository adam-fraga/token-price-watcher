package exec

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strings"
)

type Token struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func BeautifyOutput(token Token) {
	name := color.New(color.FgCyan).SprintFunc()
	symbol := color.New(color.FgYellow).SprintFunc()
	id := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("🪙 Name: %-30s | Symbol: %-10s | ID: %s\n",
		name(token.Name),
		symbol(token.Symbol),
		id(token.ID),
	)
}

func SearchTokenID(query string) {
	keyword := strings.ToLower(query)
	url := "https://api.coingecko.com/api/v3/coins/list"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lecture:", err)
		return
	}

	var tokens []Token
	if err := json.Unmarshal(body, &tokens); err != nil {
		fmt.Println("Erreur parsing JSON:", err)
		return
	}

	fmt.Printf("🔍 Résultats pour '%s': ", keyword)
	found := false
	for _, token := range tokens {
		if strings.Contains(strings.ToLower(token.Name), keyword) || strings.Contains(strings.ToLower(token.Symbol), keyword) {
			BeautifyOutput(token)
			found = true
		}
	}

	if !found {
		fmt.Println("❌ Aucun token trouvé.")
	}
}
