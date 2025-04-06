package services

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type FileConfig struct {
	Path            string `yaml:"path"`
	RemoteDirectory string `yaml:"remote_directory"`
}

func ArchiveFiles(filePaths []string, archivePath string) error {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, filePath := range filePaths {
		err := addToArchive(zipWriter, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func addToArchive(zipWriter *zip.Writer, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Si le chemin est un répertoire, parcourir les fichiers
	if info.IsDir() {
		return filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return writeToArchive(zipWriter, path, filePath)
		})
	}
	return writeToArchive(zipWriter, filePath, filepath.Dir(filePath))
}

func writeToArchive(zipWriter *zip.Writer, path, baseDir string) error {
	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		return err
	}

	writer, err := zipWriter.Create(relPath)
	if err != nil {
		return err
	}

	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	}
	return nil
}
