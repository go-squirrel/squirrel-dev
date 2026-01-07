package terminal

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type TerminalHandler interface {
	Write(data []byte) (int, error)
	Read(output []byte) (int, error)
	Resize(width, height int) error
	Close() error
}

func NewTerminalHandler(terminalType string, height, width int, newParam any) (TerminalHandler, error) {

	if terminalType == "ssh" {
		value, ok := newParam.(*ssh.Client)
		if !ok {
			return nil, fmt.Errorf("err newParam is not *ssh.Client")
		}
		sshTerminal, err := NewSShTerminal(value, height, width)
		if err != nil {
			return nil, err
		}
		return sshTerminal, nil
	}
	return nil, nil
}
