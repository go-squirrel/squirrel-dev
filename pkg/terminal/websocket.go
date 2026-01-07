package terminal

import (
	"io"
	"log"
	"sync"

	"github.com/gorilla/websocket"
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
			log.Printf("关闭终端处理器错误: %v", err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("关闭WebSocket连接错误: %v", err)
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
					log.Printf("写入WebSocket错误: %v", err)
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("读取终端输出错误: %v", err)
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
					log.Printf("读取WebSocket错误: %v", err)
				}
				return
			}

			switch msg.Type {
			case "stdin":
				if _, err := handler.Write([]byte(msg.Data)); err != nil {
					log.Printf("写入终端输入错误: %v", err)
					return
				}
			case "resize":
				if err := handler.Resize(msg.Cols, msg.Rows); err != nil {
					log.Printf("调整终端大小错误: %v", err)
				}
			default:
				log.Printf("未知的消息类型: %s", msg.Type)
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
