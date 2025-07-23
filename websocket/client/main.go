package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		var message string
		_, err := fmt.Scanln(&message)
		if err != nil {
			log.Fatal(err)
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println(err)
			return
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received message: %s\n", string(msg))
	}
}
