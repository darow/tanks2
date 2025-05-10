package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	input Input

	x float32
	y float32

	boardImage *ebiten.Image

	tiles Tiles
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT
}

func (g *Game) Update() error {
	g.input.Update()

	if g.input.xIncrease {
		g.x++
	}

	if g.input.xDecrease {
		g.x--
	}

	if g.input.yIncrease {
		g.y++
	}

	if g.input.yDecrease {
		g.y--
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		b := Bullet{
			id: g.tiles.getNextID(),
			x:  g.x + MOUSE_CHARACTER_WIDTH/2,
			y:  g.y + MOUSE_CHARACTER_WIDTH/2,
		}

		g.tiles.bullets[b.id] = b
	}

	for i, b := range g.tiles.bullets {
		b.y++
		g.tiles.bullets[i] = b
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//if g.boardImage == nil {
	//	g.boardImage = ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	//}

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(1)
	op.GeoM.Reset()
	op.GeoM.Translate(float64(g.x), float64(g.y))

	g.boardImage.Clear()
	g.boardImage.Fill(color.RGBA{0xca, 0xff, 0xcd, 0xff})

	for _, b := range g.tiles.bullets {
		vector.DrawFilledCircle(g.boardImage, b.x, b.y, float32(BULLET_RADIUS), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
	}

	g.boardImage.DrawImage(MOUSE_CHARACTER_IMAGE, op)

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}
