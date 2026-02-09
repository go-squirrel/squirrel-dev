package terminal

import (
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 读取缓冲区大小
	bufSize = 1024
)

// WSMessage WebSocket消息类型
type WSMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Cols int    `json:"cols,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

// HandleWebSocket 处理WebSocket连接和终端IO
func HandleWebSocket(conn *websocket.Conn, handler TerminalHandler) {
	defer func() {
		if err := handler.Close(); err != nil {
			zap.L().Error("failed to close terminal handler", zap.Error(err))
		}
		if err := conn.Close(); err != nil {
			zap.L().Error("failed to close websocket connection", zap.Error(err))
		}
	}()

	var wg sync.WaitGroup

	// 从终端读取数据并发送到WebSocket
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, bufSize)
		for {
			n, err := handler.Read(buf)
			if n > 0 {
				msg := WSMessage{
					Type: "stdout",
					Data: string(buf[:n]),
				}
				if err := conn.WriteJSON(msg); err != nil {
					zap.L().Error("failed to write to websocket", zap.Error(err))
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

	// 从WebSocket读取数据并写入终端
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			var msg WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					zap.L().Error("failed to read from websocket", zap.Error(err))
				}
				return
			}

			switch msg.Type {
			case "stdin":
				if _, err := handler.Write([]byte(msg.Data)); err != nil {
					zap.L().Error("failed to write to terminal", zap.Error(err))
					return
				}
			case "resize":
				if err := handler.Resize(msg.Cols, msg.Rows); err != nil {
					zap.L().Error("failed to resize terminal",
						zap.Int("cols", msg.Cols),
						zap.Int("rows", msg.Rows),
						zap.Error(err),
					)
				}
			default:
				zap.L().Warn("unknown websocket message type", zap.String("type", msg.Type))
			}
		}
	}()

	wg.Wait()
}

// WriteMessage 发送消息到WebSocket
func WriteMessage(conn *websocket.Conn, msgType string, data string) error {
	msg := WSMessage{
		Type: msgType,
		Data: data,
	}
	return conn.WriteJSON(msg)
}

// ReadMessage 从WebSocket读取消息
func ReadMessage(conn *websocket.Conn) (*WSMessage, error) {
	var msg WSMessage
	if err := conn.ReadJSON(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
