package config

//Paramètres généraux (tokens, API, RPC, clés, etc.)
import (
	"github.com/joho/godotenv"
	"log"
)

func InitData() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Erreur chargement .env")
	}
}
