package websocket

import (
	"github.com/comoc-im/message/auth"
	"github.com/comoc-im/message/signal"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	maxMessageSize = 8 * 1024 // 8MB
	authWait       = 10 * time.Second
)

func Auth(conn *websocket.Conn) (*Client, error) {
	log.Info("authentication", zap.String("host", conn.UnderlyingConn().RemoteAddr().String()))

	// message size limit
	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(authWait))

	// read
	_, rawBytes, err := conn.ReadMessage()
	if err != nil {
		log.Error("read first message error", zap.Error(err))
		return nil, err
	}

	si := auth.SignIn{}
	if err := si.Decode(&rawBytes); err != nil {
		log.Error("bad sign in message", zap.Error(err))
		return nil, err
	}

	return &Client{
		address: si.Address,
		conn:    conn,
		signal:  make(chan *signal.Signal, 16),
	}, nil
}
