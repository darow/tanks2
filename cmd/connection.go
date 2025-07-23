package main

import (
	"encoding/json"
	"log"

	"myebiten/cmd/websocket/client"
	"myebiten/cmd/websocket/server"

	"github.com/gorilla/websocket"
)

const (
	CONNECTION_MODE_OFFLINE = "offline"
	CONNECTION_MODE_SERVER  = "server"
	CONNECTION_MODE_CLIENT  = "client"
)

func (g *Game) makeSuccessConnection() {
	switch *CONNECTION_MODE {
	case CONNECTION_MODE_SERVER:
		g.server = server.New()
	case CONNECTION_MODE_CLIENT:
		g.client = client.New()
	default:
	}
	SUCCESS_CONNECTION = true
}

func (g *Game) UpdateFromClient() {
	msg, err := g.client.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	var newGame Game
	err = json.Unmarshal(msg, &newGame)
	if err != nil {
		log.Fatal(err)
	}

	g.charactersScores = newGame.charactersScores

	g.characters[0] = newGame.characters[0]
	g.characters[1].x = newGame.characters[1].x
	g.characters[1].y = newGame.characters[1].y
	g.characters[1].rotation = newGame.characters[1].rotation

	g.things = newGame.things
}

func (c *Character) SendInputToServer(ws *websocket.Conn) {
	msg, err := json.Marshal(c.input)
	if err != nil {
		log.Fatal(err)
	}

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Fatal(err)
	}
}
