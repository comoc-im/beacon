package main

import (
	"errors"
	"github.com/comoc-im/message"
	"github.com/comoc-im/message/address"
	"github.com/comoc-im/message/auth"
	"github.com/comoc-im/message/signal"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
)

var conMap map[address.Address]*websocket.Conn

func handleConnection(conn *websocket.Conn) {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("error close connection", err)
		}
	}(conn)
	for {
		mt, rawBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if len(rawBytes) == 0 {
			log.Println(errors.New("empty message"))
			continue
		}

		firstByte := rawBytes[0]
		messageType := message.MessageType(firstByte)
		switch messageType {
		case message.SignIn:
			si := auth.SignIn{}
			if err := si.Decode(&rawBytes); err != nil {
				log.Println("bad sign in message", err)
				continue
			}
			if _, ok := conMap[si.Address]; !ok {
				defer delete(conMap, si.Address)
			}
			conMap[si.Address] = conn
			continue
		case message.Signal:
			s := signal.Signal{}
			if err := s.Decode(&rawBytes); err != nil {
				log.Println("bad sign in message", err)
				continue
			}
			target, ok := conMap[s.To]
			if !ok {
				log.Println("target:", err)
				continue
			}
			err = target.WriteMessage(mt, rawBytes)
			if err != nil {
				log.Println("write:", err)
				continue
			}
		default:
			log.Println(errors.New("not valid message"))
			continue
		}
	}
}

func beacon(res http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//log.Println("req origin", r.Header["Origin"])
			return true
		},
	}

	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go handleConnection(conn)
}

func main() {
	conMap = map[address.Address]*websocket.Conn{}
	handler := cors.Default().Handler(http.HandlerFunc(beacon))
	log.Println("beacon server up and running.")
	if err := http.ListenAndServe("127.0.0.1:9999", handler); err != nil {
		panic(err)
	}
}
