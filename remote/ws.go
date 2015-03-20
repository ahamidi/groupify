package main

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

func NewWSConnection(url string, origin string) (*websocket.Conn, error) {
	log.Println("Connecting to Groupify API..")

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}
	log.Println("Connected!")
	return ws, nil
}

func SendWSMessage(conn *websocket.Conn, message interface{}) error {

	// Convert input to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send it!
	_, err = conn.Write(payload)
	if err != nil {
		return err
	}
	return nil
}

func WSMessageReceiver(conn *websocket.Conn, notifications chan interface{}) error {

	for {
		msg := make([]byte, 512)
		msgLen, err := conn.Read(msg)
		if err != nil {
			log.Println("Recieved Message:", string(msg[:msgLen]))
			log.Println("Error:", err)
		}

		// Unmarshal the message
		log.Println("Received Message:", string(msg[:msgLen]))
		notifications <- msg[:msgLen]
	}
}
