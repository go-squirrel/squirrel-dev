package terminal

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/crypto/ssh"
)

type TerminalSession struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
	mu      sync.Mutex
}

func (ts *TerminalSession) Write(data []byte) (int, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.stdin == nil {
		return 0, fmt.Errorf("终端未初始化")
	}

	return ts.stdin.Write(data)
}

func (ts *TerminalSession) Read(output []byte) (int, error) {
	if ts.stdout == nil {
		return 0, fmt.Errorf("终端未初始化")
	}

	return ts.stdout.Read(output)
}

func (ts *TerminalSession) Resize(width, height int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.session == nil {
		return fmt.Errorf("会话未初始化")
	}

	return ts.session.WindowChange(height, width)
}

func (ts *TerminalSession) Close() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.session != nil {
		ts.session.Close()
	}

	if ts.stdin != nil {
		ts.stdin.Close()
	}

	return nil
}
