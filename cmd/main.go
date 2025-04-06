package main

import (
	"flag"
	"log"

	"github.com/codiaby/backup-service/services"
)

func main() {
	// Options pour le chemin de configuration et exécution immédiate
	configPath := flag.String("C", "config/config.yaml", "Chemin du fichier de configuration YAML")
	flag.StringVar(configPath, "config", "config/config.yaml", "Chemin du fichier de configuration YAML (long format)")
	runNow := flag.Bool("run-now", false, "Exécuter immédiatement les sauvegardes sans planification")
	flag.Parse()

	// Démarrage du service principal
	err := services.StartBackupService(*configPath, *runNow)
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du service : %v", err)
	}
}
