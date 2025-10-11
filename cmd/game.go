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
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nfnt/resize"
)

const (
	STATE_GAME_ENDING_TIMER_SECONDS = 4
)

const (
	STATE_MAP_CREATING = iota
	STATE_GAME_RUNNING
	STATE_GAME_ENDING
)

type Game struct {
	stateEndingTimer *time.Timer
	state            int
	leftAlive        int
	boardSizeX       int
	boardSizeY       int
	Things           Things
	Characters       []*Character
	CharactersStash  []*Character
	CharactersScores []uint

	server     *server.Server
	client     *client.Client
	boardImage *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	for i, char := range g.Characters {
		if char == nil {
			continue
		}

		g.Characters[i].input.Update()

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

			continue
		} else if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
			if i == 1 {
				msg := g.server.ReadMessage()

				var newInput Input
				err := json.Unmarshal(msg, &newInput)
				if err != nil {
					continue
				}

				char.input = newInput

				newBullet := g.Characters[i].Update(g.Things.walls)
				if newBullet != nil {
					newBullet.ID = g.Things.getNextID()
					g.Things.Bullets[newBullet.ID] = *newBullet
				}
			} else {
				newBullet := g.Characters[i].Update(g.Things.walls)
				if newBullet != nil {
					newBullet.ID = g.Things.getNextID()
					g.Things.Bullets[newBullet.ID] = *newBullet
				}
			}
		}
	}

	if *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
		g.UpdateGameFromServer()
		return nil
	} else if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
		switch g.state {
		case STATE_MAP_CREATING:
			g.CreateMap()
			g.SendMapToClient()
			g.leftAlive = 2
			g.state = STATE_GAME_RUNNING
		case STATE_GAME_RUNNING:
			if g.leftAlive <= 1 {
				g.stateEndingTimer = time.NewTimer(STATE_GAME_ENDING_TIMER_SECONDS * time.Second)
				g.state = STATE_GAME_ENDING
			}
		case STATE_GAME_ENDING:
			select {
			case <-g.stateEndingTimer.C:
				for _, char := range g.Characters {
					if char != nil {
						g.CharactersScores[char.id]++
						break
					}
				}
				g.state = STATE_MAP_CREATING
			default:
			}
		default:
		}
	}

	for i, b := range g.Things.Bullets {
		if b.X < 0 || b.X > float64(SCREEN_SIZE_WIDTH) || b.Y < 0 || b.Y > float64(SCREEN_SIZE_HEIGHT) {
			delete(g.Things.Bullets, i)
			continue
		}

		g.Things.Bullets[i] = g.Things.ProcessBullet(b)
	}

	for bulletKey, bullet := range g.Things.Bullets {
		for charIndex, char := range g.Characters {
			if char == nil {
				continue
			}
			isCollision := g.Things.DetectBulletCharacterCollision(bullet, char)
			if isCollision {
				delete(g.Things.Bullets, bulletKey)
				g.leftAlive--

				if FEATURE_DECREASING_TANKS {
					g.Characters[charIndex].CurrentWidth--
					resizedCharacterImage := resize.Resize(g.Characters[charIndex].CurrentWidth, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
					g.Characters[charIndex].charImg = ebiten.NewImageFromImage(resizedCharacterImage)
					continue
				}

				charToStash := char
				if charIndex != char.id {
					g.Characters = append(g.Characters, char)
					lastIndex := len(g.Characters) - 1
					g.Characters[lastIndex].id = lastIndex
					charToStash = g.Characters[lastIndex]
				}

				g.Characters[charIndex] = nil
				g.CharactersStash = append(g.CharactersStash, charToStash)
			}
		}
	}

	if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
		timeToSendMap := time.Now().UnixNano()%10000 == 0

		if timeToSendMap {
			g.SendMapToClient()
		}

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
	//if g.boardImage == nil {
	//	g.boardImage = ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	//}

	g.boardImage.Clear()
	g.boardImage.Fill(color.RGBA{0xca, 0xca, 0xff, 0xff})

	for _, b := range g.Things.Bullets {
		vector.DrawFilledCircle(g.boardImage, float32(b.X), float32(b.Y), float32(BULLET_RADIUS), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
	}

	g.Things.wallsMu.RLock()
	for w := range g.Things.walls {
		width, height := float32(WALL_WIDTH), float32(WALL_HEIGHT)
		if w.Horizontal {
			width, height = WALL_HEIGHT, WALL_WIDTH
		}

		vector.DrawFilledRect(g.boardImage, float32(w.X)*WALL_HEIGHT, float32(w.Y)*WALL_HEIGHT, width, height, color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)

		if DEBUG_MODE {
			vector.DrawFilledCircle(g.boardImage, float32(w.X)*WALL_HEIGHT, float32(w.Y)*WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
			vector.DrawFilledCircle(g.boardImage, float32(w.X)*WALL_HEIGHT+WALL_WIDTH, float32(w.Y)*WALL_HEIGHT+WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
		}
	}
	g.Things.wallsMu.RUnlock()

	for _, c := range g.Characters {
		if c == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(c.CurrentWidth)/2, -float64(c.CurrentWidth)/2)
		op.GeoM.Rotate(c.Rotation)
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(c.X, c.Y)
		g.boardImage.DrawImage(c.charImg, op)

		vector.DrawFilledCircle(g.boardImage, float32(c.X), float32(c.Y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)

		if DEBUG_MODE {
			charCorners := c.getCorners()
			for _, corner := range charCorners {
				vector.DrawFilledCircle(g.boardImage, float32(corner.x), float32(corner.y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
			}
		}
	}

	scoreX := g.boardSizeX * WALL_HEIGHT / 2
	scoreY := g.boardSizeY*WALL_HEIGHT + 40
	text.Draw(g.boardImage, fmt.Sprintf("Player 1: %d        Player 2: %d", g.CharactersScores[0], g.CharactersScores[1]), REGULAR_FONT, scoreX, scoreY, color.Black)

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}
