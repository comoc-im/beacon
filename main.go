package main

import (
	"encoding/json"
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
	From    string      `json:"from"`
	To      string      `json:"to"`
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

var conMap map[string]*websocket.Conn

func beacon(res http.ResponseWriter, req *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	query := req.URL.Query()
	username := query.Get("username")

	delete(conMap, username)
	conMap[username] = conn

	defer conn.Close()
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			log.Printf("recv: %s", message)

			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("parse:", err)
				continue
			}

			var targetConn *websocket.Conn
			if msg.Type == Heartbeat {
				targetConn = conn
			} else {
				target, ok := conMap[msg.To]
				if !ok {
					log.Println("target:", err)
					continue
				}
				targetConn = target
			}

			err = targetConn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}()
}

func main() {
	conMap = map[string]*websocket.Conn{}
	handler := cors.Default().Handler(http.HandlerFunc(beacon))
	log.Println("beacon server up and running.")
	if err := http.ListenAndServe(":9999", handler); err != nil {
		panic(err)
	}
}
