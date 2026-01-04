package ssh

import gossh "golang.org/x/crypto/ssh"

type Client struct {
	Client *gossh.Client
}
