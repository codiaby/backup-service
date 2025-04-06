package main

import (
	"backup-service/services"
	"flag"
	"log"
)

func main() {
	// Option pour le chemin du fichier de configuration
	configPath := flag.String("C", "config/config.yaml", "Chemin du fichier de configuration YAML")
	flag.StringVar(configPath, "config", "config/config.yaml", "Chemin du fichier de configuration YAML (long format)")
	flag.Parse()

	// Démarrage du service principal
	err := services.StartBackupService(*configPath)
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du service : %v", err)
	}
}
