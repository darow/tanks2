package main

import (
	"image/color"

	"myebiten/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	COLOR_BLACK      = color.RGBA{0x0f, 0x0f, 0x0f, 0xff}
	COLOR_BACKGROUND = color.RGBA{0xca, 0xca, 0xff, 0xff}
)

type Hitbox interface {
	Hit(this *models.GameObject, hb Hitbox, other *models.GameObject)
}

type Sprite interface {
	Draw(drawingArea *models.DrawingArea, gameObject models.GameObject)
}

type RectangleHitbox struct {
	W, H float64
}

func (rectHB RectangleHitbox) Hit(this *models.GameObject, hb Hitbox, other *models.GameObject) {

}

type RectangleSprite struct {
	W, H float64
}

func (rectangleSprite RectangleSprite) Draw(drawingArea *models.DrawingArea, gameObject models.GameObject) {
	width, height := rectangleSprite.W, rectangleSprite.H
	if gameObject.Rotation == 0.0 {
		width, height = height, width
	}
	topLeftCorner := models.Vector2D{X: gameObject.Position.X - width/2, Y: gameObject.Position.Y - height/2}

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

type ImageSprite struct {
	*ebiten.Image
}

func (imageSprite ImageSprite) Draw(drawingArea *models.DrawingArea, gameObject models.GameObject) {
	op := &ebiten.DrawImageOptions{}

	w := imageSprite.Image.Bounds().Max.X - imageSprite.Image.Bounds().Min.X
	h := imageSprite.Image.Bounds().Max.Y - imageSprite.Image.Bounds().Min.Y

	sc := drawingArea.Scale
	offX := drawingArea.Offset.X
	offY := drawingArea.Offset.Y

	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(gameObject.Rotation)
	op.GeoM.Scale(sc, sc)
	op.GeoM.Translate(gameObject.Position.X*sc+offX, gameObject.Position.Y*sc+offY)

	image := drawingArea.BoardImage

	image.DrawImage(imageSprite.Image, op)
}
