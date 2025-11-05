package models

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Bullet struct {
	GameObject
	R float64
}

func (b Bullet) Draw(drawingArea *DrawingArea) {
	sc := drawingArea.Scale
	offX := drawingArea.Offset.X
	offY := drawingArea.Offset.Y

	x := float32(b.Position.X)*float32(sc) + float32(offX)
	y := float32(b.Position.Y)*float32(sc) + float32(offY)
	r := float32(b.R) * float32(sc)

	image := drawingArea.BoardImage

	vector.DrawFilledCircle(image, x, y, r, color.Black, false)
}
