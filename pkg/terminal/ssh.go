package terminal

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SShTerminal struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
	mu      sync.Mutex
}

func NewSShTerminal(sshClient *ssh.Client, height, width int) (*SShTerminal, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}

	if err := session.RequestPty(term, height, width, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("请求PTY失败: %v", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, err
	}

	if err := session.Shell(); err != nil {
		session.Close()
		return nil, fmt.Errorf("启动shell失败: %v", err)
	}

	return &SShTerminal{
		session: session,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}, nil
}

func (ts *SShTerminal) Write(data []byte) (int, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.stdin == nil {
		return 0, fmt.Errorf("终端未初始化")
	}

	return ts.stdin.Write(data)
}

func (ts *SShTerminal) Read(output []byte) (int, error) {
	if ts.stdout == nil {
		return 0, fmt.Errorf("终端未初始化")
	}

	return ts.stdout.Read(output)
}

func (ts *SShTerminal) Resize(width, height int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.session == nil {
		return fmt.Errorf("会话未初始化")
	}

	return ts.session.WindowChange(height, width)
}

func (ts *SShTerminal) Close() error {
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
