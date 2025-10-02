package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

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
		g.client = client.New(*ADDRESS)
		go g.ReceiveMapUpdates()
	default:
	}
	SUCCESS_CONNECTION = true
}

func (g *Game) UpdateGameFromServer() {
	msg := g.client.ReadMessage()

	var newGame Game
	err := json.Unmarshal(msg, &newGame)
	if err != nil {
		log.Println(err)
		return
	}

	if len(newGame.Characters) == 0 && len(newGame.CharactersStash) == 0 {
		return
	}

	g.CharactersScores = newGame.CharactersScores

	if len(g.Characters) == 0 {
		g.Characters = append(g.Characters, newGame.Characters...)
	} else {
		g.Characters[0].Copy(newGame.Characters[0])
		g.Characters[1].Copy(newGame.Characters[1])
	}

	g.Things.Bullets = newGame.Things.Bullets
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

func (g *Game) ReceiveMapUpdates() {
	for {
		_, message, err := g.client.MapUpdateConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
		}

		fmt.Printf("Received map message: %s\n", message)

		var wallsDTO WallsDTO
		err = json.Unmarshal(message, &wallsDTO)
		if err != nil {
			log.Fatal(err)
		}

		newWalls := make(map[Wall]struct{})

		for _, w := range wallsDTO.Walls {
			newWalls[w] = struct{}{}
		}

		g.Things.wallsMu.Lock()
		g.Things.walls = newWalls
		g.Things.wallsMu.Unlock()
	}
}
