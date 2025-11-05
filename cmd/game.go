package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"time"

	"myebiten/cmd/websocket/client"
	"myebiten/cmd/websocket/server"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	STATE_GAME_ENDING_TIMER_SECONDS = 4
	ITEM_SPAWN_INTERVAL             = 5
)

const (
	STATE_MAP_CREATING = iota
	STATE_GAME_RUNNING
	STATE_GAME_ENDING
)

type Game struct {
	stateEndingTimer *time.Timer
	itemSpawnTicker  *time.Ticker

	state     int
	leftAlive int

	Maze             [][]MazeNode
	Bullets          []*Bullet
	Walls            []Wall
	Characters       []*Character
	CharactersScores []uint

	mainArea *DrawingArea
	UIArea1  *DrawingArea
	UIArea2  *DrawingArea

	server     *server.Server
	client     *client.Client
	boardImage *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	for i, char := range g.Characters {
		if !char.Active {
			continue
		}

		char.input.Update()

		if i == 0 && *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
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

			// continue
		} else if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
			if i == 1 {
				msg := g.server.ReadMessage()

				var newInput Input
				err := json.Unmarshal(msg, &newInput)
				if err != nil {
					continue
				}

				char.input = newInput

				char.ProcessInput()
			} else {
				char.ProcessInput()
			}
		} else {
			char.ProcessInput()
		}

		char.Move()
		g.DetectCharacterToWallCollision(char)
	}

	for _, bullet := range g.Bullets {
		bullet.Move()
	}

	for _, bullet := range g.Bullets {
		if !bullet.Active {
			continue
		}

		g.DetectBulletToWallCollision(bullet)

		for _, char := range g.Characters {
			if !char.Active {
				continue
			}

			if false { //g.DetectBulletToCharacterCollision(bullet, char) {
				bullet.Active = false
				char.Active = false
				g.leftAlive--

				// if FEATURE_DECREASING_TANKS {
				// 	g.Characters[charIndex].CurrentWidth--
				// 	resizedCharacterImage := resize.Resize(g.Characters[charIndex].CurrentWidth, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
				// 	g.Characters[charIndex].charImg = ebiten.NewImageFromImage(resizedCharacterImage)
				// 	continue
				// }
			}
		}
	}

	if *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
		g.UpdateGameFromServer()
		return nil
	} else /* if *CONNECTION_MODE == CONNECTION_MODE_SERVER */ {
		switch g.state {
		case STATE_MAP_CREATING:
			g.Reset()
			g.SetupLevel()
			g.itemSpawnTicker = time.NewTicker(ITEM_SPAWN_INTERVAL * time.Second)
			// g.SendMapToClient()
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
				g.state = STATE_MAP_CREATING
			default:
			}
		default:
		}
	}

	// if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
	// 	msg, err := json.Marshal(g)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	err = g.server.WriteThingsMessage(msg)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.boardImage.Clear()
	g.boardImage.Fill(COLOR_BACKGROUND)

	for _, wall := range g.Walls {
		if wall.Active {
			wall.Draw(g.mainArea, wall.GameObject)
		}
	}

	for _, character := range g.Characters {
		if character.Active {
			character.Draw(g.mainArea, character.GameObject)
		}
	}

	for _, bullet := range g.Bullets {
		if bullet.Active {
			bullet.Draw(g.mainArea, bullet.GameObject)
		}
	}
	// if DEBUG_MODE {
	// 	vector.DrawFilledCircle(g.boardImage, float32(w.X)*WALL_HEIGHT, float32(w.Y)*WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
	// 	vector.DrawFilledCircle(g.boardImage, float32(w.X)*WALL_HEIGHT+WALL_WIDTH, float32(w.Y)*WALL_HEIGHT+WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
	// }
	// if DEBUG_MODE {
	// 	charCorners := c.getCorners()
	// 	for _, corner := range charCorners {
	// 		vector.DrawFilledCircle(g.boardImage, float32(corner.x), float32(corner.y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
	// 	}
	// }

	text.Draw(g.boardImage, fmt.Sprintf("Player 1: %d        Player 2: %d", g.CharactersScores[0], g.CharactersScores[1]), REGULAR_FONT, 0, 0, color.Black)

	// for i := range g.boardSizeY {
	// 	for j := range g.boardSizeX {
	// 		cx := j*WALL_HEIGHT + WALL_HEIGHT/2
	// 		cy := i*WALL_HEIGHT + WALL_HEIGHT/2

	// 		wh := float32(WALL_HEIGHT)
	// 		ww := float32(WALL_WIDTH)

	// 		vector.DrawFilledCircle(g.boardImage, float32(cx), float32(cy)-(wh-ww)/2, 1, color.RGBA{0xFF, 0x0, 0x0, 0x0}, false)
	// 		vector.DrawFilledCircle(g.boardImage, float32(cx)-(wh-ww)/2, float32(cy), 1, color.RGBA{0xFF, 0x0, 0x0, 0x0}, false)
	// 		vector.DrawFilledCircle(g.boardImage, float32(cx), float32(cy), 1, color.White, false)
	// 	}
	// }

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}
