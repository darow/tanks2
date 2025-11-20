package game

import (
	"encoding/json"
	"errors"
	"image/color"
	"log"
	"time"

	"myebiten/internal/models"
	"myebiten/internal/websocket/client"
	"myebiten/internal/websocket/server"

	"github.com/hajimehoshi/ebiten/v2"
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

var wallsToCheck []*models.Wall = make([]*models.Wall, 12)
var (
	COLOR_BLACK = color.RGBA{0x0f, 0x0f, 0x0f, 0xff}
)

type Game struct {
	stateEndingTimer *time.Timer
	itemSpawnTicker  *time.Ticker

	state     int
	leftAlive int

	Maze             [][]MazeNode
	Bullets          []*models.Bullet
	Walls            []models.Wall
	Characters       []*models.Character
	CharactersScores []uint

	server   *server.Server
	client   *client.Client
	connMode string

	scoreUITexts []models.UIText
	scenes       map[int]*models.Scene
	activeScene  *models.Scene
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	if g.connMode == CONNECTION_MODE_CLIENT {
		char := g.Characters[0]

		char.Input.Update()

		msg, err := json.Marshal(char.Input)
		if err != nil {
			log.Fatal(err)
		}

		err = g.client.WriteMessage(msg)
		if err != nil {
			log.Fatal(err)
		}

		if char.Input.Shoot {
			char.Input.Shoot = false
		}

		g.UpdateGameFromServer()

		return nil
	}

	switch g.state {
	case STATE_MAZE_CREATING:
		g.Reset()
		g.itemSpawnTicker = time.NewTicker(ITEM_SPAWN_INTERVAL * time.Second)
		h, w, walls := g.SetupLevel()
		if g.connMode != CONNECTION_MODE_OFFLINE {
			g.SendMazeToClient(h, w, walls)
		}
		g.leftAlive = 2
		g.state = STATE_GAME_RUNNING

		g.SanityCheck()

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
				if char.IsActive() {
					g.updateScores(char.ID)
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
		if !char.IsActive() {
			continue
		}

		if i == 1 && g.connMode == CONNECTION_MODE_SERVER {
			// process client's character's input
			msg := g.server.ReadMessage()

			var input models.Input
			err := json.Unmarshal(msg, &input)
			if err != nil {
				continue
			}

			char.Input = input
		} else {
			char.Input.Update()
		}

		char.ProcessInput()

		char.Move()

		g.DetectCharacterToWallCollision(char)
	}

	for _, bullet := range g.Bullets {
		if !bullet.IsActive() {
			continue
		}

		bullet.Move()

		g.DetectBulletToWallCollision(bullet)

		for _, char := range g.Characters {
			if !char.IsActive() {
				continue
			}

			if char.DetectBulletToCharacterCollision(bullet) {
				bullet.SetActive(false)
				char.SetActive(false)
				g.leftAlive--
			}
		}
	}

	if g.connMode == CONNECTION_MODE_SERVER {
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
	image := g.activeScene.Draw()

	screen.Clear()
	screen.DrawImage(image, &ebiten.DrawImageOptions{})
}

func CreateGame(bullets []*models.Bullet, characters []*models.Character, scenes map[int]*models.Scene, scoreUIs []models.UIText) *Game {
	return &Game{
		Bullets:          bullets,
		Characters:       characters,
		CharactersScores: []uint{0, 0, 0, 0},
		scenes:           scenes,
		scoreUITexts:     scoreUIs,
	}
}
