package websocket

import (
	"github.com/comoc-im/message/signal"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

func Read(client *Client) {
	defer func() {
		client.close()
	}()

	// keepalive
	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		log.Info("heartbeat PONG")
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// read
	for {
		_, rawBytes, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("read message error", zap.Error(err))
			}
			break
		}

		s := signal.Signal{}
		if err := s.Decode(&rawBytes); err != nil {
			log.Error("bad signal message", zap.Error(err))
			continue
		}

		hub.transfer(&s)
	}
}
