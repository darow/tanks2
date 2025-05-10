package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	input Input

	character Character
	tiles     Tiles

	boardImage *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT
}

func (g *Game) Update() error {
	g.input.Update()

	if g.input.rotateRight {
		g.character.rotation += CHARACTER_ROTATION_SPEED
	}

	if g.input.rotateLeft {
		g.character.rotation -= CHARACTER_ROTATION_SPEED
	}

	if g.input.moveForward {
		sin, cos := math.Sincos(g.character.rotation)
		g.character.x += float32(cos) * CHARACTER_SPEED
		g.character.y += float32(sin) * CHARACTER_SPEED
	}

	if g.input.moveBackward {
		sin, cos := math.Sincos(g.character.rotation)
		g.character.x -= float32(cos) * CHARACTER_SPEED
		g.character.y -= float32(sin) * CHARACTER_SPEED
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		sin, cos := math.Sincos(g.character.rotation)
		x := g.character.x - float32(cos)*CHARACTER_WIDTH/2
		y := g.character.y - float32(sin)*CHARACTER_WIDTH/2

		b := Bullet{
			id:       g.tiles.getNextID(),
			x:        x,
			y:        y,
			rotation: g.character.rotation,
		}

		g.tiles.bullets[b.id] = b
	}

	for i, b := range g.tiles.bullets {
		sin, cos := math.Sincos(b.rotation)
		b.x -= float32(cos) * BULLET_SPEED
		b.y -= float32(sin) * BULLET_SPEED

		if !(b.x > 0 && b.x < SCREEN_SIZE_WIDTH && b.y > 0 && b.y < SCREEN_SIZE_HEIGHT) {
			delete(g.tiles.bullets, i)
			continue
		}

		// vertical wall collisions
		if int(b.x)%WALL_HEIGHT <= WALL_WIDTH {
			WallToCollide := Wall{
				x:          uint16(math.Round(float64(b.x / WALL_HEIGHT))),
				y:          uint16(math.Floor(float64(b.y / WALL_HEIGHT))),
				horizontal: false,
			}

			if _, ok := g.tiles.walls[WallToCollide]; ok {
				var newRotation float64
				if b.rotation < math.Pi {
					newRotation = math.Pi - b.rotation
				} else {
					newRotation = b.rotation - math.Pi
				}

				b.rotation = newRotation
			}
		}
		// horizontal wall collisions
		if int(b.y)%WALL_HEIGHT <= WALL_WIDTH {
			WallToCollide := Wall{
				x:          uint16(math.Round(float64(b.x / WALL_HEIGHT))),
				y:          uint16(math.Floor(float64(b.y / WALL_HEIGHT))),
				horizontal: true,
			}

			if _, ok := g.tiles.walls[WallToCollide]; ok {

				b.rotation = -b.rotation
			}
		}

		g.tiles.bullets[i] = b

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//if g.boardImage == nil {
	//	g.boardImage = ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	//}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(-CHARACTER_WIDTH/2, -CHARACTER_WIDTH/2)
	op.GeoM.Rotate(g.character.rotation)
	op.GeoM.Translate(float64(g.character.x), float64(g.character.y))

	g.boardImage.Clear()
	g.boardImage.Fill(color.RGBA{0xca, 0xca, 0xff, 0xff})

	for _, b := range g.tiles.bullets {
		vector.DrawFilledCircle(g.boardImage, b.x, b.y, float32(BULLET_RADIUS), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
	}

	for w, _ := range g.tiles.walls {
		width, height := float32(WALL_WIDTH), float32(WALL_HEIGHT)
		if w.horizontal {
			width, height = WALL_HEIGHT, WALL_WIDTH
		}

		vector.DrawFilledRect(g.boardImage, float32(w.x)*WALL_HEIGHT, float32(w.y)*WALL_HEIGHT, width, height, color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)

		vector.DrawFilledCircle(g.boardImage, float32(w.x)*WALL_HEIGHT, float32(w.y)*WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)
		vector.DrawFilledCircle(g.boardImage, float32(w.x)*WALL_HEIGHT+WALL_WIDTH, float32(w.y)*WALL_HEIGHT+WALL_HEIGHT, 2, color.RGBA{0x00, 0xff, 0xff, 0xff}, false)

	}

	g.boardImage.DrawImage(CHARACTER_IMAGE, op)

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}
