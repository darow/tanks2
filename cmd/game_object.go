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
}

type Sprite interface {
	Draw(image *ebiten.Image, center Vector2D, rotation float64)
}

type RectangleHitbox struct {
	w, h float64
}

type CircleHitbox struct {
	r float64
}

type RectangleSprite struct {
	w, h float64
}

func (rectangleSprite RectangleSprite) Draw(image *ebiten.Image, center Vector2D, rotation float64) {
	width, height := rectangleSprite.w, rectangleSprite.h
	if rotation == 0.0 {
		width, height = height, width
	}
	topLeftCorner := Vector2D{center.x - width/2, center.y - height/2}

	x := float32(topLeftCorner.x)*float32(DRAWING_SCALE) + float32(DRAWING_OFFSET_X)
	y := float32(topLeftCorner.y)*float32(DRAWING_SCALE) + float32(DRAWING_OFFSET_Y)
	w := float32(width) * float32(DRAWING_SCALE)
	h := float32(height) * float32(DRAWING_SCALE)

	// if takes the coordinates in pixels, WHY THE FUCK ARE THEY FLOAT32???????????
	vector.DrawFilledRect(image, x, y, w, h, color.Black, false)
}

type BallSprite struct {
	r float64
}

func (ballSprite BallSprite) Draw(image *ebiten.Image, center Vector2D, rotation float64) {
	x := float32(center.x)*float32(DRAWING_SCALE) + float32(DRAWING_OFFSET_X)
	y := float32(center.y)*float32(DRAWING_SCALE) + float32(DRAWING_OFFSET_Y)
	r := float32(ballSprite.r) * float32(DRAWING_SCALE)

	vector.DrawFilledCircle(image, x, y, r, color.Black, false)
}

type ImageSprite struct {
	*ebiten.Image
}

func (imageSprite ImageSprite) Draw(image *ebiten.Image, center Vector2D, rotation float64) {
	op := &ebiten.DrawImageOptions{}

	w := imageSprite.Image.Bounds().Max.X - imageSprite.Image.Bounds().Min.X
	h := imageSprite.Image.Bounds().Max.Y - imageSprite.Image.Bounds().Min.Y

	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(rotation)
	op.GeoM.Scale(float64(DRAWING_SCALE), float64(DRAWING_SCALE))
	op.GeoM.Translate(center.x*DRAWING_SCALE+float64(DRAWING_OFFSET_X), center.y*DRAWING_SCALE+float64(DRAWING_OFFSET_Y))

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

// It won't allow me to pass the pointer to the object
func (dw DefaultWeapon) Shoot(origin Vector2D, rotation float64) {
	for _, bullet := range dw.clip {
		if !bullet.active {
			bullet.position.x = origin.x
			bullet.position.y = origin.y

			bullet.rotation = rotation

			sin, cos := math.Sincos(rotation)
			bullet.speed.x = cos * BULLET_SPEED
			bullet.speed.y = sin * BULLET_SPEED

			bullet.active = true
			break
		}
	}

	dw.Discharge()
}

// Same here
func (dw DefaultWeapon) Discharge() {

}

type GameObject struct {
	id       int
	active   bool
	position Vector2D
	rotation float64
	speed    Vector2D
}

func (gameObject *GameObject) Move() {
	gameObject.position.x += gameObject.speed.x
	gameObject.position.y += gameObject.speed.y
}
