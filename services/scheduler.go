package services

import (
	"log"
	"sync"
	"time"

	"os"

	"github.com/codiaby/backup-service/services"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Databases []services.DatabaseConfig `yaml:"databases"`
	Files     []services.FileConfig     `yaml:"files"`
	Backup    struct {
		Directory      string `yaml:"directory"`
		RetentionDays  int    `yaml:"retention_days"`
		MaxConcurrency int    `yaml:"max_concurrency"`
	} `yaml:"backup"`
	Server   services.ServerConfig `yaml:"server"`
	Schedule string                `yaml:"schedule"`
}

func StartBackupService(configPath string) error {
	config, err := loadConfig(configPath)
	if err != nil {
		return err
	}

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

func performBackup(config *Config) {
	maxConcurrency := config.Backup.MaxConcurrency
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	// Gestion des bases de données
	for _, db := range config.Databases {
		wg.Add(1)
		go func(db services.DatabaseConfig) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			timestamp := time.Now().Format("20060102_150405")
			backupFile := config.Backup.Directory + db.Name + "_" + timestamp + ".sql"

			if err := services.BackupDatabase(db, backupFile); err != nil {
				log.Printf("Erreur lors de la sauvegarde de %s : %v", db.Name, err)
				return
			}

			if err := services.UploadToServer(services.ServerConfig{
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

	// Gestion des fichiers
	for _, file := range config.Files {
		wg.Add(1)
		go func(file services.FileConfig) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			timestamp := time.Now().Format("20060102_150405")
			archivePath := config.Backup.Directory + "archive_" + timestamp + ".zip"
			if err := services.ArchiveFiles([]string{file.Path}, archivePath); err != nil {
				log.Printf("Erreur lors de l'archivage des fichiers : %v", err)
				return
			}

			if err := services.UploadToServer(services.ServerConfig{
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

	wg.Wait()
	log.Println("Toutes les sauvegardes et les transferts sont terminés.")

	// Nettoyage des fichiers anciens
	log.Println("Nettoyage des fichiers anciens...")
	err := services.CleanupOldBackups(config.Backup.Directory, config.Backup.RetentionDays)
	if err != nil {
		log.Printf("Erreur lors du nettoyage des fichiers anciens : %v", err)
	}
}
