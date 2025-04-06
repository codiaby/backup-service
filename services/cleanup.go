package services

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func CleanupOldBackups(backupDir string, retentionDays int) error {
	// Calcul de la limite de temps
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	// Parcours du répertoire des sauvegardes
	err := filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Suppression des fichiers anciens
		if !info.IsDir() && info.ModTime().Before(cutoff) {
			log.Printf("Suppression du fichier ancien : %s", path)
			err := os.Remove(path)
			if err != nil {
				log.Printf("Erreur lors de la suppression de %s : %v", path, err)
			}
		}
		return nil
	})
	return err
}
