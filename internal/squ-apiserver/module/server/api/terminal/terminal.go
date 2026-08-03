package terminal

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const bufferSize = 1024

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Cols int    `json:"cols,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

type SSH struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	mu      sync.Mutex
}

func NewSSH(client *ssh.Client, height, width int) (*SSH, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}
	modes := ssh.TerminalModes{
		ssh.ECHO: 1, ssh.TTY_OP_ISPEED: 14400, ssh.TTY_OP_OSPEED: 14400,
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
	if _, err := session.StderrPipe(); err != nil {
		session.Close()
		return nil, err
	}
	if err := session.Shell(); err != nil {
		session.Close()
		return nil, fmt.Errorf("启动shell失败: %v", err)
	}
	return &SSH{session: session, stdin: stdin, stdout: stdout}, nil
}

func (s *SSH) Write(value []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stdin == nil {
		return 0, fmt.Errorf("终端未初始化")
	}
	return s.stdin.Write(value)
}

func (s *SSH) Read(value []byte) (int, error) {
	if s.stdout == nil {
		return 0, fmt.Errorf("终端未初始化")
	}
	return s.stdout.Read(value)
}

func (s *SSH) Resize(width, height int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.session == nil {
		return fmt.Errorf("会话未初始化")
	}
	return s.session.WindowChange(height, width)
}

func (s *SSH) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.session != nil {
		_ = s.session.Close()
	}
	if s.stdin != nil {
		_ = s.stdin.Close()
	}
	return nil
}

func Bridge(conn *websocket.Conn, handler *SSH) {
	defer func() {
		if err := handler.Close(); err != nil {
			zap.L().Error("failed to close terminal handler", zap.Error(err))
		}
		if err := conn.Close(); err != nil {
			zap.L().Error("failed to close terminal websocket", zap.Error(err))
		}
	}()
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		defer wait.Done()
		buffer := make([]byte, bufferSize)
		for {
			n, err := handler.Read(buffer)
			if n > 0 {
				if writeErr := conn.WriteJSON(Message{Type: "stdout", Data: string(buffer[:n])}); writeErr != nil {
					zap.L().Error("failed to write terminal output to websocket", zap.Error(writeErr))
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					zap.L().Error("failed to read from terminal", zap.Error(err))
				}
				return
			}
		}
	}()
	go func() {
		defer wait.Done()
		for {
			var message Message
			if err := conn.ReadJSON(&message); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					zap.L().Error("failed to read from terminal websocket", zap.Error(err))
				}
				return
			}
			switch message.Type {
			case "stdin":
				if _, err := handler.Write([]byte(message.Data)); err != nil {
					zap.L().Error("failed to write to terminal", zap.Error(err))
					return
				}
			case "resize":
				if err := handler.Resize(message.Cols, message.Rows); err != nil {
					zap.L().Error("failed to resize terminal",
						zap.Int("cols", message.Cols),
						zap.Int("rows", message.Rows),
						zap.Error(err),
					)
				}
			default:
				zap.L().Warn("unknown terminal websocket message type", zap.String("type", message.Type))
			}
		}
	}()
	wait.Wait()
}

func WriteMessage(conn *websocket.Conn, messageType, data string) error {
	return conn.WriteJSON(Message{Type: messageType, Data: data})
}
