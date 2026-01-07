package terminal

import (
	"fmt"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type SSHSession struct {
	conn   *ssh.Client
	client *ssh.Client
	mu     sync.Mutex
}

func NewSSHSession(sshClient *ssh.Client) (*SSHSession, error) {
	return &SSHSession{
		conn:   sshClient,
		client: sshClient,
	}, nil
}

func (s *SSHSession) StartTerminal(width, height int) (*TerminalSession, error) {
	session, err := s.client.NewSession()
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

	return &TerminalSession{
		session: session,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}, nil
}

func (s *SSHSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return s.client.Close()
	}

	return nil
}

// HandleWebSocketWithSSH 处理WebSocket连接并使用SSH终端
func HandleWebSocketWithSSH(conn *websocket.Conn, sshClient *ssh.Client) error {
	session, err := NewSSHSession(sshClient)
	if err != nil {
		return err
	}
	defer session.Close()

	terminalSession, err := session.StartTerminal(80, 24)
	if err != nil {
		return err
	}

	HandleWebSocket(conn, terminalSession)
	return nil
}
