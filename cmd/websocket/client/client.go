package client

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
)

type MessageStore struct {
	sync.Mutex
	message []byte
}

type Client struct {
	charInputConn    *websocket.Conn
	thingsUpdateConn *websocket.Conn
	MapUpdateConn    *websocket.Conn
	msgStore         *MessageStore
}

func New(hostAddress string) *Client {
	address := fmt.Sprintf("ws://%s/ws1", hostAddress)
	charInputConn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	address = fmt.Sprintf("ws://%s/ws2", hostAddress)
	thingsUpdateConn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	address = fmt.Sprintf("ws://%s/ws3", hostAddress)
	mapUpdateConn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	c := &Client{
		charInputConn:    charInputConn,
		thingsUpdateConn: thingsUpdateConn,
		MapUpdateConn:    mapUpdateConn,
		msgStore:         &MessageStore{},
	}
	go c.ReceiveUpdates()

	return c
}

func (c *Client) ReceiveUpdates() {
	for {
		_, message, err := c.thingsUpdateConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
		}

		// fmt.Printf("Received message: %s\n", message)

		c.msgStore.Lock()
		c.msgStore.message = message
		c.msgStore.Unlock()
	}
}

func (c *Client) ReadMessage() []byte {
	c.msgStore.Lock()
	message := c.msgStore.message
	c.msgStore.Unlock()

	return message
}

func (c *Client) WriteMessage(message []byte) error {
	if err := c.charInputConn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}
