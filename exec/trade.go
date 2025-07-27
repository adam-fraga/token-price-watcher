// Exécution logique : sell command, logging, alertes
package exec

import (
	"fmt"
	"os/exec"
)

func ExecuteSell(tokenID string, price float64) {
	fmt.Printf("[🚀] %s atteint %.2f USD. Exécution de la vente...\n", tokenID, price)
	// Appelle un script de vente ou intègre Web3 ici
	cmd := exec.Command("sh", "sell.sh", tokenID, fmt.Sprintf("%.2f", price))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Erreur exécution commande : %v\n", err)
	}
}
