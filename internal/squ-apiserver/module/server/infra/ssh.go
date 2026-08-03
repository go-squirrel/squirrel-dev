package infra

import (
	"context"

	"squirrel-dev/internal/squ-apiserver/module/server/domain"
	sshClient "squirrel-dev/pkg/ssh"
)

type SSHTester struct{}

func NewSSHTester() *SSHTester { return &SSHTester{} }

func (SSHTester) Test(_ context.Context, server domain.Server) error {
	client, err := sshClient.NewSsh(machine(server, server.Hostname))
	if err != nil {
		return err
	}
	client.Close()
	return nil
}

func NewSSHClient(server domain.Server) (*sshClient.Client, error) {
	return sshClient.NewSsh(machine(server, "test"))
}

func machine(server domain.Server, name string) *sshClient.Machine {
	password := ""
	if server.SSHPassword != nil {
		password = *server.SSHPassword
	}
	privateKey := ""
	if server.SSHPrivateKey != nil {
		privateKey = *server.SSHPrivateKey
	}
	return &sshClient.Machine{
		Name: name, IpAddress: server.IPAddress, User: server.SSHUsername,
		Password: password, Port: server.SSHPort, PrivateKey: privateKey, Type: server.AuthType,
	}
}
