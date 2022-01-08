package websocket

import (
	"github.com/gorilla/websocket"
	"time"
)

func Write(client *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-client.signal:
			if !ok {
				return
			}

			log.Info("send signal")
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))

			err := client.conn.WriteMessage(websocket.BinaryMessage, message.Encode())
			if err != nil {
				return
			}
		case <-ticker.C:
			log.Info("heartbeat PING")
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
