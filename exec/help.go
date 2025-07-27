package exec

import "fmt"

func PrintHelp() {
	fmt.Println(`🛠️  Utilisation :

  -search <mot-clé>     → Recherche les tokens liés au nom ou symbole donné
  -freq <minutes>       → Fréquence d'extraction du prix (ex: -freq 10 pour toutes les 10 min)
  -limit <prix>         → Définir une limite de prix pour vendre (ex: -limit 34.5)
  -help                 → Affiche cette aide

📌 Exemples :
  go run cmd/main.go -search avax
  go run cmd/main.go -search avax -freq 5 -limit 25.3
`)
}
