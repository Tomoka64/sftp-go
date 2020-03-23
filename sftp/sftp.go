package sftp

import (
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTP interface {
	Upload(filename string, data []byte) error
	Exist(filename string) (bool, error)
	Close() error
}

type Client struct {
	cli *sftp.Client
}

func NewClient(sshClient *ssh.Client) (*Client, error) {
	cli, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return &Client{cli}, nil
}

func (cl *Client) Upload(filename string, data []byte) error {
	f, err := cl.cli.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func (cl *Client) Exist(filename string) bool {
	f, err := cl.cli.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !f.IsDir()
}

func (cl *Client) Close() error {
	return cl.cli.Close()
}
