package main

import (
	"encoding/hex"
	"fmt"
	"github.com/comoc-im/message"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type MessageType string

const (
	Description MessageType = "description"
	Candidate               = "candidate"
	Heartbeat               = "heartbeat"
)

type Message struct {
	From    message.Address `json:"from"`
	To      message.Address `json:"to"`
	Type    MessageType     `json:"type"`
	Payload string          `json:"payload"`
}

var conMap map[message.Address]*websocket.Conn

func beacon(res http.ResponseWriter, req *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			log.Println("req origin", r.Header["Origin"])
			return true
		},
	}

	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	query := req.URL.Query()
	username := query.Get("username")

	log.Println(username, "connecting")
	data, err := hex.DecodeString(username)
	if err != nil {
		panic(err)
	}
	fmt.Printf("% x", data)
	add := message.Address{}
	copy(add[:], data)
	delete(conMap, add)
	conMap[add] = conn

	go func() {
		defer conn.Close()
		for {
			mt, rawBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Println("recv:", mt, rawBytes)

			if len(rawBytes) == 1 && rawBytes[0] == message.PING {
				err = conn.WriteMessage(mt, []byte{message.PONG})
				if err != nil {
					log.Println("heartbeat:", err)
					break
				}
				log.Println("heartbeat")
				continue
			}

			//var msg Message
			msg := message.Signal{}
			if err := msg.Decode(&rawBytes); err != nil {
				log.Println("parse:", err)
				continue
			}

			var targetConn *websocket.Conn
			target, ok := conMap[msg.To]
			if !ok {
				log.Println("target:", err)
				continue
			}
			targetConn = target
			err = targetConn.WriteMessage(mt, rawBytes)
			if err != nil {
				log.Println("write:", err)
				break
			}

		}
	}()
}

func main() {
	conMap = map[message.Address]*websocket.Conn{}
	handler := cors.Default().Handler(http.HandlerFunc(beacon))
	log.Println("beacon server up and running.")
	if err := http.ListenAndServe("127.0.0.1:9999", handler); err != nil {
		panic(err)
	}
}
