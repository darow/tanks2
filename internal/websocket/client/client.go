package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	mapMsgStore      *MessageStore
	playerID         int
}

func New(hostAddress string, playerID int) *Client {
	address := fmt.Sprintf("ws://%s/ws1?player_id=%d", hostAddress, playerID)
	charInputConn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	address = fmt.Sprintf("ws://%s/ws2?player_id=%d", hostAddress, playerID)
	thingsUpdateConn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	address = fmt.Sprintf("ws://%s/ws3?player_id=%d", hostAddress, playerID)
	mapUpdateConn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	c := &Client{
		charInputConn:    charInputConn,
		thingsUpdateConn: thingsUpdateConn,
		MapUpdateConn:    mapUpdateConn,
		msgStore:         &MessageStore{},
		mapMsgStore:      &MessageStore{},
		playerID:         playerID,
	}
	go c.ReceiveUpdates()
	go c.ReceiveMapUpdates()

	return c
}

func GetPlayersCount(hostAddress string) (int, error) {
	address := fmt.Sprintf("http://%s/players_count", hostAddress)
	response, err := http.Get(address)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, errors.New(response.Status)
	}

	var playersCount int
	if err := json.NewDecoder(response.Body).Decode(&playersCount); err != nil {
		return 0, err
	}

	return playersCount, nil
}

func (c *Client) ReceiveUpdates() {
	for {
		_, message, err := c.thingsUpdateConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
		}

		c.msgStore.Lock()
		c.msgStore.message = message
		c.msgStore.Unlock()
	}
}

func (c *Client) ReceiveMapUpdates() {
	for {
		_, message, err := c.MapUpdateConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
		}

		c.mapMsgStore.Lock()
		c.mapMsgStore.message = message
		c.mapMsgStore.Unlock()
	}
}

func (c *Client) ReadMessage() []byte {
	c.msgStore.Lock()
	message := c.msgStore.message
	c.msgStore.Unlock()

	return message
}

func (c *Client) ReadMapMessage() []byte {
	c.mapMsgStore.Lock()
	message := c.mapMsgStore.message
	c.mapMsgStore.message = nil
	c.mapMsgStore.Unlock()

	return message
}

func (c *Client) WriteMessage(message []byte) error {
	if err := c.charInputConn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetPlayerID() int {
	return c.playerID
}
