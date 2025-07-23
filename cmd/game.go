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
	boardSizeX       int
	boardSizeY       int
	things           Things
	characters       []*Character
	charactersStash  []*Character
	charactersScores []uint

	server     *server.Server
	client     *client.Client
	boardImage *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	if !SUCCESS_CONNECTION {
		g.makeSuccessConnection()
	}

	for i, char := range g.characters {
		if char == nil {
			continue
		}
		g.characters[i].input.Update()
		if i == 1 && *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
			msg, err := json.Marshal(char.input)
			if err != nil {
				log.Fatal(err)
			}

			err = g.client.WriteMessage(msg)
			if err != nil {
				log.Fatal(err)
			}

			continue
		} else if i == 1 && *CONNECTION_MODE == CONNECTION_MODE_SERVER {
			msg, err := g.server.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			var newInput Input
			err = json.Unmarshal(msg, &newInput)
			if err != nil {
				log.Fatal(err)
			}

			char.input = newInput
			continue
		}

		newBullet := g.characters[i].Update(g.things.walls)
		if newBullet != nil {
			newBullet.id = g.things.getNextID()
			g.things.bullets[newBullet.id] = *newBullet
		}
	}

	if *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
		g.UpdateFromClient()
		return nil
	}

	switch g.state {
	case STATE_MAP_CREATING:
		g.CreateMap()
		g.state = STATE_GAME_RUNNING
	case STATE_GAME_ENDING:
		select {
		case <-g.stateEndingTimer.C:
			g.state = STATE_MAP_CREATING
		default:
		}
	default:
	}

	for i, b := range g.things.bullets {
		if b.x < 0 || b.x > float64(SCREEN_SIZE_WIDTH) || b.y < 0 || b.y > float64(SCREEN_SIZE_HEIGHT) {
			delete(g.things.bullets, i)
			continue
		}

		g.things.bullets[i] = g.things.ProcessBullet(b)
	}

	for bulletKey, bullet := range g.things.bullets {
		for charIndex, char := range g.characters {
			if char == nil {
				continue
			}
			isCollision := g.things.DetectBulletCharacterCollision(bullet, char)
			if isCollision {
				delete(g.things.bullets, bulletKey)

				if FEATURE_DECREASING_TANKS {
					g.characters[charIndex].currentWidth--
					resizedCharacterImage := resize.Resize(g.characters[charIndex].currentWidth, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
					g.characters[charIndex].charImg = ebiten.NewImageFromImage(resizedCharacterImage)
					continue
				}
				switch g.state {
				case STATE_GAME_RUNNING:
					for _, otherChar := range g.characters {
						if otherChar.id != char.id {
							g.charactersScores[otherChar.id]++
						}
					}
					g.state = STATE_GAME_ENDING
					g.stateEndingTimer = time.NewTimer(STATE_GAME_ENDING_TIMER_SECONDS * time.Second)
				case STATE_GAME_ENDING:
					{
						g.charactersScores[char.id]--
					}
				}

				charToStash := char
				if charIndex != char.id {
					g.characters = append(g.characters, char)
					lastIndex := len(g.characters) - 1
					g.characters[lastIndex].id = lastIndex
					charToStash = g.characters[lastIndex]
				}

				g.characters[charIndex] = nil
				g.charactersStash = append(g.charactersStash, charToStash)
			}
		}
	}

	if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
		msg, err := json.Marshal(g)
		if err != nil {
			log.Fatal(err)
		}
		err = g.server.WriteMessage(msg)
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

	for _, b := range g.things.bullets {
		vector.DrawFilledCircle(g.boardImage, float32(b.x), float32(b.y), float32(BULLET_RADIUS), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
	}

	for w, _ := range g.things.walls {
		width, height := float32(WALL_WIDTH), float32(WALL_HEIGHT)
		if w.horizontal {
			width, height = WALL_HEIGHT, WALL_WIDTH
		}

		vector.DrawFilledRect(g.boardImage, float32(w.x)*WALL_HEIGHT, float32(w.y)*WALL_HEIGHT, width, height, color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)

		if DEBUG_MODE {
			vector.DrawFilledCircle(g.boardImage, float32(w.x)*WALL_HEIGHT, float32(w.y)*WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
			vector.DrawFilledCircle(g.boardImage, float32(w.x)*WALL_HEIGHT+WALL_WIDTH, float32(w.y)*WALL_HEIGHT+WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
		}
	}

	for _, c := range g.characters {
		if c == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(c.currentWidth)/2, -float64(c.currentWidth)/2)
		op.GeoM.Rotate(c.rotation)
		op.GeoM.Translate(c.x, c.y)
		g.boardImage.DrawImage(c.charImg, op)

		vector.DrawFilledCircle(g.boardImage, float32(c.x), float32(c.y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)

		if DEBUG_MODE {
			charCorners := c.getCorners()
			for _, corner := range charCorners {
				vector.DrawFilledCircle(g.boardImage, float32(corner.x), float32(corner.y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
			}
		}
	}

	scoreX := g.boardSizeX * WALL_HEIGHT / 2
	scoreY := g.boardSizeY*WALL_HEIGHT + 40
	text.Draw(g.boardImage, fmt.Sprintf("%d %d", g.charactersScores[0], g.charactersScores[1]), REGULAR_FONT, scoreX, scoreY, color.Black)

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}
