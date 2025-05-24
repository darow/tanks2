package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nfnt/resize"
)

type Game struct {
	characters []*Character
	tiles      Tiles

	boardImage *ebiten.Image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT
}

func (g *Game) Update() error {
	for i := range g.characters {
		g.characters[i].input.Update()
		newBullet := g.characters[i].Update(g.tiles.walls)
		if newBullet != nil {
			newBullet.id = g.tiles.getNextID()
			g.tiles.bullets[newBullet.id] = *newBullet
		}
	}

	for i, b := range g.tiles.bullets {
		if !(b.x > 0 && b.x < SCREEN_SIZE_WIDTH && b.y > 0 && b.y < SCREEN_SIZE_HEIGHT) {
			delete(g.tiles.bullets, i)
			continue
		}

		g.tiles.bullets[i] = g.tiles.ProcessBullet(b)
	}

	for bulletKey, bullet := range g.tiles.bullets {
		for charIndex, char := range g.characters {
			isCollision := g.tiles.DetectBulletCharacterCollision(bullet, char)
			if isCollision {
				delete(g.tiles.bullets, bulletKey)

				if FEATURE_DECREASING_TANKS {
					g.characters[charIndex].currentWidth--
					resizedCharacterImage := resize.Resize(g.characters[charIndex].currentWidth, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
					g.characters[charIndex].charImg = ebiten.NewImageFromImage(resizedCharacterImage)
				}
			}
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

	for _, b := range g.tiles.bullets {
		vector.DrawFilledCircle(g.boardImage, float32(b.x), float32(b.y), float32(BULLET_RADIUS), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
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

	for _, c := range g.characters {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(c.currentWidth)/2, -float64(c.currentWidth)/2)
		op.GeoM.Rotate(c.rotation)
		op.GeoM.Translate(c.x, c.y)
		g.boardImage.DrawImage(c.charImg, op)

		vector.DrawFilledCircle(g.boardImage, float32(c.x), float32(c.y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
		charCorners := c.getCorners()
		for _, corner := range charCorners {
			vector.DrawFilledCircle(g.boardImage, float32(corner.x), float32(corner.y), float32(1), color.RGBA{0x0f, 0x0f, 0x0f, 0xff}, false)
		}
	}

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}
