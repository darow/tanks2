package models

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Vector2D struct {
	X, Y float64
}

func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

type GameObject struct {
	ID       int
	Active   bool
	Position Vector2D
	Rotation float64
	Speed    Vector2D
}

func (gameObject *GameObject) IsActive() bool {
	return gameObject.Active
}

func (gameObject *GameObject) SetActive(b bool) {
	gameObject.Active = b
}

func (gameObject *GameObject) Move() {
	gameObject.Position.X += gameObject.Speed.X
	gameObject.Position.Y += gameObject.Speed.Y
}

func (gameObject *GameObject) MoveBack() {
	gameObject.Position.X -= gameObject.Speed.X
	gameObject.Position.Y -= gameObject.Speed.Y
}

type RectangleSprite struct {
	W, H float64
}

func (rectangleSprite RectangleSprite) Draw(centerX, centerY, rotation float64, drawingArea *DrawingArea) {
	width, height := rectangleSprite.W, rectangleSprite.H

	if rotation == 0.0 {
		width, height = height, width
	}
	topLeftCorner := Vector2D{X: centerX - width/2, Y: centerY - height/2}

	sc := drawingArea.Scale
	offX := drawingArea.Offset.X
	offY := drawingArea.Offset.Y

	x := float32(topLeftCorner.X)*float32(sc) + float32(offX)
	y := float32(topLeftCorner.Y)*float32(sc) + float32(offY)
	w := float32(width) * float32(sc)
	h := float32(height) * float32(sc)

	image := drawingArea.BoardImage

	vector.DrawFilledRect(image, x, y, w, h, color.Black, false)
}

type CircleSprite struct {
	R float64
}

func (circleSprite CircleSprite) Draw(centerX, centerY float64, drawingArea *DrawingArea) {
	sc := drawingArea.Scale
	offX := drawingArea.Offset.X
	offY := drawingArea.Offset.Y

	x := float32(centerX)*float32(sc) + float32(offX)
	y := float32(centerY)*float32(sc) + float32(offY)
	r := float32(circleSprite.R) * float32(sc)

	image := drawingArea.BoardImage

	vector.DrawFilledCircle(image, x, y, r, color.Black, false)
}

type ImageSprite struct {
	*ebiten.Image
}

func (imageSprite ImageSprite) Draw(centerX, centerY, rotation float64, drawingArea *DrawingArea) {
	op := &ebiten.DrawImageOptions{}

	w := imageSprite.Image.Bounds().Max.X - imageSprite.Image.Bounds().Min.X
	h := imageSprite.Image.Bounds().Max.Y - imageSprite.Image.Bounds().Min.Y

	sc := drawingArea.Scale
	offX := drawingArea.Offset.X
	offY := drawingArea.Offset.Y

	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(rotation)
	op.GeoM.Scale(sc, sc)
	op.GeoM.Translate(centerX*sc+offX, centerY*sc+offY)

	image := drawingArea.BoardImage

	image.DrawImage(imageSprite.Image, op)
}

type RectangleHitbox struct {
	H, W float64
}
