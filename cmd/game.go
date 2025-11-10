package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"log"
	"time"

	"myebiten/cmd/websocket/client"
	"myebiten/cmd/websocket/server"
	"myebiten/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	STATE_GAME_ENDING_TIMER_SECONDS = 4
	ITEM_SPAWN_INTERVAL             = 5
)

const (
	STATE_MAZE_CREATING = iota
	STATE_GAME_RUNNING
	STATE_GAME_ENDING
)

var wallsToCheck []*Wall = make([]*Wall, 12)

type Game struct {
	stateEndingTimer *time.Timer
	itemSpawnTicker  *time.Ticker

	state     int
	leftAlive int

	Maze             [][]MazeNode
	Bullets          models.Pool
	Walls            []Wall
	Characters       []*Character
	CharactersScores []uint

	mainArea *models.DrawingArea
	UIArea1  *models.DrawingArea
	UIArea2  *models.DrawingArea

	server     *server.Server
	client     *client.Client
	boardImage *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	if *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
		char := g.Characters[0]

		char.input.Update()

		msg, err := json.Marshal(char.input)
		if err != nil {
			log.Fatal(err)
		}

		err = g.client.WriteMessage(msg)
		if err != nil {
			log.Fatal(err)
		}

		if char.input.Shoot {
			char.input.Shoot = false
		}

		g.UpdateGameFromServer()

		return nil
	}

	switch g.state {
	case STATE_MAZE_CREATING:
		g.Reset()
		g.itemSpawnTicker = time.NewTicker(ITEM_SPAWN_INTERVAL * time.Second)
		h, w, walls := g.SetupLevel()
		if *CONNECTION_MODE != CONNECTION_MODE_OFFLINE {
			g.SendMazeToClient(h, w, walls)
		}
		g.leftAlive = 2
		g.state = STATE_GAME_RUNNING

	case STATE_GAME_RUNNING:
		select {
		case <-g.itemSpawnTicker.C:
			g.SpawnItem()
		default:
			if g.leftAlive <= 1 {
				g.stateEndingTimer = time.NewTimer(STATE_GAME_ENDING_TIMER_SECONDS * time.Second)
				g.state = STATE_GAME_ENDING
			}
		}

	case STATE_GAME_ENDING:
		select {
		case <-g.stateEndingTimer.C:
			for _, char := range g.Characters {
				if char.Active {
					g.CharactersScores[char.ID]++
					break
				}
			}
			g.state = STATE_MAZE_CREATING
		default:
		}

	default:
		return errors.New("invalid state")
	}

	for i, char := range g.Characters {
		if !char.Active {
			continue
		}

		if i == 1 && *CONNECTION_MODE == CONNECTION_MODE_SERVER {
			// process client's character's input
			msg := g.server.ReadMessage()

			var input Input
			err := json.Unmarshal(msg, &input)
			if err != nil {
				continue
			}

			char.input = input
		} else {
			char.input.Update()
		}

		char.ProcessInput()

		char.Move()

		g.DetectCharacterToWallCollision(char)
	}

	for _, bullet := range g.Bullets.Elements() {
		if !bullet.Active {
			continue
		}

		bullet.Move()

		g.DetectBulletToWallCollision(bullet)

		for _, char := range g.Characters {
			if !char.Active {
				continue
			}

			if g.DetectBulletToCharacterCollision(bullet, char) {
				bullet.Active = false
				char.Active = false
				g.leftAlive--
			}
		}
	}

	if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
		msg, err := json.Marshal(g)
		if err != nil {
			log.Fatal(err)
		}
		err = g.server.WriteThingsMessage(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.boardImage.Clear()
	g.boardImage.Fill(COLOR_BACKGROUND)

	mazeArea := g.mainArea.Children[0]

	for _, wall := range g.Walls {
		if wall.Active {
			wall.Draw(mazeArea, wall.GameObject)
		}
	}

	for _, character := range g.Characters {
		if character.Active {
			character.Draw(mazeArea, character.GameObject)
		}
	}

	for _, bullet := range g.Bullets.Elements() {
		if bullet.Active {
			bullet.Draw(mazeArea)
		}
	}

	for i, score := range g.CharactersScores {
		DrawScore(g.UIArea2.Children[i], fmt.Sprintf("Player %d", i+1), score)
	}

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}

func DrawScore(drawingArea *models.DrawingArea, name string, score uint) {
	text.Draw(drawingArea.BoardImage, fmt.Sprintf("%s: %d", name, score), REGULAR_FONT, int(drawingArea.Offset.X), int(drawingArea.Offset.Y), color.Black)
}
