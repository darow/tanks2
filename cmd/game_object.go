package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	COLOR_BLACK      = color.RGBA{0x0f, 0x0f, 0x0f, 0xff}
	COLOR_BACKGROUND = color.RGBA{0xca, 0xca, 0xff, 0xff}
)

type Hitbox interface {
	Hit(this *GameObject, hb Hitbox, other *GameObject)
}

type Sprite interface {
	Draw(drawingArea *DrawingArea, gameObject GameObject)
}

type RectangleHitbox struct {
	w, h float64
}

func (rectHB RectangleHitbox) Hit(this *GameObject, hb Hitbox, other *GameObject) {

}

type CircleHitbox struct {
	r float64
}

func (circleHB CircleHitbox) Hit(this *GameObject, hb Hitbox, other *GameObject) {

}

type RectangleSprite struct {
	W, H float64
}

func (rectangleSprite RectangleSprite) Draw(drawingArea *DrawingArea, gameObject GameObject) {
	width, height := rectangleSprite.W, rectangleSprite.H
	if gameObject.Rotation == 0.0 {
		width, height = height, width
	}
	topLeftCorner := Vector2D{gameObject.Position.x - width/2, gameObject.Position.y - height/2}

	sc := drawingArea.Scale
	offX := drawingArea.Offset.x
	offY := drawingArea.Offset.y

	x := float32(topLeftCorner.x)*float32(sc) + float32(offX)
	y := float32(topLeftCorner.y)*float32(sc) + float32(offY)
	w := float32(width) * float32(sc)
	h := float32(height) * float32(sc)

	image := drawingArea.boardImage

	// if takes the coordinates in pixels, WHY THE FUCK ARE THEY FLOAT32???????????
	vector.DrawFilledRect(image, x, y, w, h, color.Black, false)
}

type BallSprite struct {
	R float64
}

func (ballSprite BallSprite) Draw(drawingArea *DrawingArea, gameObject GameObject) {
	sc := drawingArea.Scale
	offX := drawingArea.Offset.x
	offY := drawingArea.Offset.y

	x := float32(gameObject.Position.x)*float32(sc) + float32(offX)
	y := float32(gameObject.Position.y)*float32(sc) + float32(offY)
	r := float32(ballSprite.R) * float32(sc)

	image := drawingArea.boardImage

	vector.DrawFilledCircle(image, x, y, r, color.Black, false)
}

type ImageSprite struct {
	*ebiten.Image
}

func (imageSprite ImageSprite) Draw(drawingArea *DrawingArea, gameObject GameObject) {
	op := &ebiten.DrawImageOptions{}

	w := imageSprite.Image.Bounds().Max.X - imageSprite.Image.Bounds().Min.X
	h := imageSprite.Image.Bounds().Max.Y - imageSprite.Image.Bounds().Min.Y

	sc := drawingArea.Scale
	offX := drawingArea.Offset.x
	offY := drawingArea.Offset.y

	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(gameObject.Rotation)
	op.GeoM.Scale(sc, sc)
	op.GeoM.Translate(gameObject.Position.x*sc+offX, gameObject.Position.y*sc+offY)

	image := drawingArea.boardImage

	image.DrawImage(imageSprite.Image, op)
}

type Weapon interface {
	Shoot(origin Vector2D, rotation float64)
	Discharge()
}

type DefaultWeapon struct {
	clip     []*Bullet
	cooldown int
}

func (dw *DefaultWeapon) Shoot(origin Vector2D, rotation float64) {
	for _, bullet := range dw.clip {
		if !bullet.Active {
			bullet.Position.x = origin.x
			bullet.Position.y = origin.y

			bullet.Rotation = rotation

			sin, cos := math.Sincos(rotation)
			bullet.Speed.x = cos * BULLET_SPEED
			bullet.Speed.y = sin * BULLET_SPEED

			bullet.Active = true
			break
		}
	}

	dw.Discharge()
}

func (dw *DefaultWeapon) Discharge() {

}

type GameObject struct {
	ID       int
	Active   bool
	Position Vector2D
	Rotation float64
	Speed    Vector2D
}

func (gameObject *GameObject) Move() {
	gameObject.Position.x += gameObject.Speed.x
	gameObject.Position.y += gameObject.Speed.y
}
