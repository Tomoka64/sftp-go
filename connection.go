package main

import (
	"context"
	"io"

	"github.com/Tomoka64/sftp-go-sample/config"

	"golang.org/x/crypto/ssh"
)

func newConnectionThroughProxy(ctx context.Context) (*ssh.Client, error) {
	proxyClient, err := startProxy(ctx)
	if err != nil {
		return nil, err
	}

	targetClient, err := getTargetConfig()
	if err != nil {
		return nil, err
	}

	proxyConn, err := proxyClient.Dial("tcp", config.TargetServer().Addr)
	if err != nil {
		return nil, err
	}
	go closeConnection(ctx, proxyConn)

	conn, chans, reqs, err := ssh.NewClientConn(proxyConn, config.TargetServer().Addr, targetClient)
	if err != nil {
		return nil, err
	}
	go closeConnection(ctx, conn)

	return ssh.NewClient(conn, chans, reqs), nil
}

func getTargetConfig() (*ssh.ClientConfig, error) {
	signer, err := ssh.ParsePrivateKeyWithPassphrase(
		[]byte(config.TargetServer().PassPhrase),
		[]byte(config.TargetServer().PrivateKey),
	)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: config.TargetServer().Addr,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

func startProxy(ctx context.Context) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKeyWithPassphrase(
		[]byte(config.Proxy().PrivateKey),
		[]byte(config.Proxy().PassPhrase),
	)
	if err != nil {
		return nil, err
	}
	proxyConfig := &ssh.ClientConfig{
		User: config.Proxy().Addr,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	proxyClient, err := ssh.Dial("tcp", config.Proxy().Addr, proxyConfig)
	if err != nil {
		return nil, err
	}
	go closeConnection(ctx, proxyClient.Conn)

	return proxyClient, nil
}

func closeConnection(ctx context.Context, conn io.Closer) {
	defer conn.Close()
	<-ctx.Done()
}
