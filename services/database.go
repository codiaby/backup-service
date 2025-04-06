package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type DatabaseConfig struct {
	Name            string `yaml:"name"`
	Type            string `yaml:"type"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	RemoteDirectory string `yaml:"remote_directory"`
}

func BackupDatabase(db DatabaseConfig, backupFile string) error {
	var cmd *exec.Cmd
	if db.Type == "mysql" {
		cmd = exec.Command("mysqldump", "-u", db.User, "-p"+db.Password, db.Name)
	} else if db.Type == "postgresql" {
		cmd = exec.Command("pg_dump", "-U", db.User, "-d", db.Name)
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
