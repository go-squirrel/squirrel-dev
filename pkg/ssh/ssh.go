package ssh

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type Machine struct {
	Name       string `yaml:"name"`
	IpAddress  string `yaml:"ip_address"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Port       int    `yaml:"port"`
	PrivateKey string `yaml:"private_key"`
	PassPhrase string
	PublicKey  string `yaml:"public_key"`
	Type       string `yaml:"type"`
}

type Client struct {
	Client *gossh.Client
	Ctx    context.Context
	Cancel context.CancelFunc
	Name   string
}

func NewSsh(machine *Machine) (s *Client, err error) {

	config := &gossh.ClientConfig{
		User:            machine.User,
		Timeout:         5 * time.Second,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}

	if machine.Type == "password" {
		config.Auth = []gossh.AuthMethod{
			gossh.Password(machine.Password),
		}
	} else {
		signer, err := makePrivateKeySigner(machine.PrivateKey, machine.PassPhrase)
		if err != nil {
			return nil, err
		}
		config.Auth = []gossh.AuthMethod{gossh.PublicKeys(signer)}
	}

	hostport := net.JoinHostPort(machine.IpAddress, fmt.Sprintf("%d", machine.Port))
	proto := "tcp"
	if strings.Contains(machine.IpAddress, ":") {
		proto = "tcp6"
	}
	client, err := gossh.Dial(proto, hostport, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	s = &Client{Client: client, Ctx: ctx, Cancel: cancel, Name: machine.Name}
	return s, err
}

func (c *Client) Close() {
	c.Client.Close()
}

func makePrivateKeySigner(privateKey string, passPhrase string) (gossh.Signer, error) {
	privateKeyByte := []byte(privateKey)
	passPhraseByte := []byte(passPhrase)

	if len(passPhraseByte) != 0 {
		return gossh.ParsePrivateKeyWithPassphrase(privateKeyByte, passPhraseByte)
	}
	return gossh.ParsePrivateKey(privateKeyByte)
}
