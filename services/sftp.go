package services

import (
	"os"
	"path/filepath"

	"github.com/codiaby/backup-service/shared"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func UploadToServer(server shared.ServerConfig, localFile string) error {
	configSSH := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", server.Address, configSSH)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	srcFile, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(filepath.Join(server.RemoteDirectory, filepath.Base(localFile)))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	return err
}
