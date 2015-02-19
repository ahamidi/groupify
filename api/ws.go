package main

import (
	"net/http"
	"log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { 
		return true 
	},
}

func messageSender(conn *websocket.Conn, messages *chan []byte) {
	log.Println("Message sender spawned")
	for m := range *messages {
		log.Println("Got message: ", m)
		if err := conn.WriteMessage(1, m); err != nil {
			log.Println(err)
		}

	}
}

func messageReceiver(conn *websocket.Conn, notifications *chan []byte) {
	log.Println("Message receiver spawned")
	for {
		messageType, m, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.Println("Received message: ", messageType, string(m))
		*notifications <- m
	}
}

func WSRemoteHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

	go messageSender(conn, context.messages)
	go messageReceiver(conn, context.notifications)
}
