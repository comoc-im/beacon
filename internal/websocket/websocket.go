package websocket

import (
	"github.com/comoc-im/beacon/internal/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//logger.Println("req origin", r.Header["Origin"])
			return true
		},
	}
	log = logger.GetLogger("WebSocket")
)

func beacon(res http.ResponseWriter, req *http.Request) {
	log.Info("new connection", zap.String("host", req.Host))
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	client, err := Auth(conn)
	if err != nil {
		return
	}

	hub.add(client)
	go Read(client)
	go Write(client)
}

func Start(addr string) {
	handler := http.HandlerFunc(beacon)
	log.Info("WebSocket server running: ", zap.String("address", addr))
	if err := http.ListenAndServe(addr, handler); err != nil {
		panic(err)
	}
}
