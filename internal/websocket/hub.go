package websocket

import (
	"github.com/comoc-im/message/address"
	"github.com/comoc-im/message/signal"
	"go.uber.org/zap"
)

var hub = Hub{
	clients: make(map[address.Address]*Client),
}

type Hub struct {
	clients map[address.Address]*Client
}

func (h *Hub) add(client *Client) {
	log.Info("Hub add")
	h.clients[client.address] = client
}

func (h *Hub) delete(client *Client) {
	log.Info("Hub delete")
	delete(h.clients, client.address)
}

func (h *Hub) transfer(s *signal.Signal) {
	log.Info("Hub transfer",
		zap.String("from", string(s.From)),
		zap.String("to", string(s.To)),
	)
	target, ok := h.clients[s.To]
	if ok && target != nil {
		target.signal <- s
	}
}
