package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/codiaby/backup-service/shared"
)

func BackupDatabase(db shared.DatabaseConfig, backupFile string) error {
	var cmd *exec.Cmd
	if db.Type == "mysql" {
		// Commande MySQL avec l'adresse du serveur
		cmd = exec.Command("mysqldump", "-h", db.Address, "-u", db.User, "-p"+db.Password, db.Name)
	} else if db.Type == "postgresql" {
		// Commande PostgreSQL avec l'adresse du serveur
		cmd = exec.Command("pg_dump", "-h", db.Address, "-U", db.User, "-d", db.Name)
		os.Setenv("PGPASSWORD", db.Password)
	} else {
		return fmt.Errorf("Type de base de données non pris en charge : %s", db.Type)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	file, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	err = cmd.Run()
	if err != nil {
		return err
	}

	_, err = file.WriteString(out.String())
	return err
}
