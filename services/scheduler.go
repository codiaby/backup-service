package services

import (
	"log"
	"sync"
	"time"

	"os"

	"github.com/codiaby/backup-service/shared"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Databases []shared.DatabaseConfig `yaml:"databases"` // Optionnel
	Files     []shared.FileConfig     `yaml:"files"`     // Optionnel
	Backup    struct {
		Directory      string `yaml:"directory"`
		RetentionDays  int    `yaml:"retention_days"`
		MaxConcurrency int    `yaml:"max_concurrency"`
		EnableCleanup  bool   `yaml:"enable_cleanup"`
	} `yaml:"backup"`
	Server   shared.ServerConfig `yaml:"server"`
	Schedule string              `yaml:"schedule"`
}

// Démarre le service de sauvegarde avec ou sans planification
func StartBackupService(configPath string, runNow bool) error {
	config, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	// Exécuter immédiatement si le drapeau --run-now est utilisé
	if runNow {
		log.Println("Exécution immédiate des sauvegardes...")
		performBackup(config)
		return nil
	}

	// Planification avec cron
	c := cron.New()
	_, err = c.AddFunc(config.Schedule, func() {
		performBackup(config)
	})
	if err != nil {
		return err
	}

	c.Start()
	log.Printf("Service de sauvegarde démarré avec l'horaire : %s", config.Schedule)
	select {}
}

// Charger la configuration depuis un fichier YAML
func loadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Effectue les sauvegardes et les transferts
func performBackup(config *Config) {
	backupDir := config.Backup.Directory
	if backupDir == "" {
		log.Println("Le champ `directory` est vide, utilisation du répertoire par défaut : ./backup/")
		backupDir = "./backup/"
	}

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Erreur lors de la création du répertoire de sauvegarde : %v", err)
			return
		}
	}

	if len(config.Databases) == 0 && len(config.Files) == 0 {
		log.Println("Aucune base de données ou fichier configuré pour la sauvegarde. Fin du processus.")
		return
	}

	maxConcurrency := config.Backup.MaxConcurrency
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	if len(config.Databases) > 0 {
		log.Println("Sauvegarde des bases de données...")
		for _, db := range config.Databases {
			wg.Add(1)
			go func(db shared.DatabaseConfig) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				timestamp := time.Now().Format("20060102_150405")
				backupFile := backupDir + db.Name + "_" + timestamp + ".sql"

				if err := BackupDatabase(db, backupFile); err != nil {
					log.Printf("Erreur lors de la sauvegarde de %s (adresse : %s) : %v", db.Name, db.Address, err)
					return
				}

				if config.Server.Address == "" || config.Server.User == "" || config.Server.Password == "" {
					log.Printf("Configuration SSH distante absente. Sauvegarde conservée en local : %s", backupFile)
					return
				}

				if err := UploadToServer(shared.ServerConfig{
					Address:         config.Server.Address,
					User:            config.Server.User,
					Password:        config.Server.Password,
					RemoteDirectory: db.RemoteDirectory,
				}, backupFile); err != nil {
					log.Printf("Erreur lors de l'envoi vers le serveur pour %s : %v", db.Name, err)
				} else {
					log.Printf("Sauvegarde et envoi réussis pour la base de données : %s", db.Name)
				}
			}(db)
		}
	} else {
		log.Println("Aucune base de données configurée. Sauvegarde ignorée.")
	}

	// Archiver et envoyer les fichiers locaux en parallèle avec des goroutines (si configurés)
	if len(config.Files) > 0 {
		log.Println("Archivage des fichiers locaux...")
		for _, file := range config.Files {
			wg.Add(1)
			go func(file shared.FileConfig) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				timestamp := time.Now().Format("20060102_150405")
				archivePath := backupDir + "archive_" + timestamp + ".zip"

				if err := ArchiveFiles([]string{file.Path}, archivePath); err != nil {
					log.Printf("Erreur lors de l'archivage des fichiers : %v", err)
					return
				}

				// Vérification des paramètres SSH
				if config.Server.Address == "" || config.Server.User == "" || config.Server.Password == "" {
					log.Printf("Configuration SSH distante absente. Archive conservée en local : %s", archivePath)
					return
				}

				// Envoi vers le serveur distant
				if err := UploadToServer(shared.ServerConfig{
					Address:         config.Server.Address,
					User:            config.Server.User,
					Password:        config.Server.Password,
					RemoteDirectory: file.RemoteDirectory,
				}, archivePath); err != nil {
					log.Printf("Erreur lors de l'envoi de l'archive vers le serveur : %v", err)
				} else {
					log.Printf("Archivage et envoi réussis : %s", archivePath)
				}
			}(file)
		}
	} else {
		log.Println("Aucun fichier configuré pour l'archivage. Archivage ignoré.")
	}

	// Attendre que toutes les tâches soient terminées
	wg.Wait()

	// Nettoyage uniquement s'il y a eu des sauvegardes ou archivages
	if len(config.Databases) > 0 || len(config.Files) > 0 {
		if config.Backup.EnableCleanup {
			log.Println("Nettoyage des fichiers anciens...")
			err := CleanupOldBackups(backupDir, config.Backup.RetentionDays)
			if err != nil {
				log.Printf("Erreur lors du nettoyage des fichiers anciens : %v", err)
			}
		} else {
			log.Println("Nettoyage des fichiers anciens ignoré (désactivé dans la configuration).")
		}
	} else {
		log.Println("Aucune sauvegarde ou fichier archivé. Nettoyage ignoré.")
	}
}
