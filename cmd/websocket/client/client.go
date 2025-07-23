package client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func New() *Client {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{conn}
}

func (c *Client) Test() {
	for {
		var message string
		_, err := fmt.Scanln(&message)
		if err != nil {
			log.Fatal(err)
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println(err)
			return
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received message: %s\n", string(msg))
	}
}

func (c *Client) ReadMessage() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Received message: %s\n", message)

	return message, err
}

func (c *Client) WriteMessage(message []byte) error {
	if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}
