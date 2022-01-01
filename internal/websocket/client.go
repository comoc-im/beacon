package websocket

import (
	"github.com/comoc-im/message/address"
	"github.com/comoc-im/message/signal"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 8) / 10
)

type Client struct {
	address address.Address
	conn    *websocket.Conn
	signal  chan *signal.Signal
}

func (client *Client) close() {
	log.Info("closing client", zap.String("client address", string(client.address)))
	_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
	hub.delete(client)
	close(client.signal)
}
