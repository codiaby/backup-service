package services

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ServerConfig struct {
	Address         string
	User            string
	Password        string
	RemoteDirectory string
}

func UploadToServer(server ServerConfig, localFile string) error {
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

	dstFile, err := client.Create(server.RemoteDirectory + "/" + filepath.Base(localFile))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
