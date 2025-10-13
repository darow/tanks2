package main

import (
	"image/color"

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
	Draw(image *ebiten.Image, center Point, rotation float64)
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

func (rectangleSprite RectangleSprite) Draw(image *ebiten.Image, center Point, rotation float64) {
	w, h := rectangleSprite.w, rectangleSprite.h
	if rotation == 0.0 {
		w, h = h, w
	}
	topLeftCorner := Point{center.x - w/2, center.y - h/2}

	// vector.StrokeLine(image, float32(topLeftCorner.x), float32(topLeftCorner.y), float32(topLeftCorner.x)+float32(w), float32(topLeftCorner.y)+float32(h), 1, color.White, false)
	// vector.StrokeLine(image, float32(topLeftCorner.x)+float32(w), float32(topLeftCorner.y), float32(topLeftCorner.x), float32(topLeftCorner.y)+float32(h), 1, color.White, false)

	// if takes the coordinates in pixels, WHY THE FUCK ARE THEY FLOAT32???????????
	vector.DrawFilledRect(image, float32(topLeftCorner.x), float32(topLeftCorner.y), float32(w), float32(h), color.Black, false)
}

type BallSprite struct {
	r float64
}

func (ballSprite BallSprite) Draw(image *ebiten.Image, center Point, rotation float64) {
	vector.DrawFilledCircle(image, float32(center.x), float32(center.y), float32(ballSprite.r), color.Black, false)
}

type ImageSprite struct {
	*ebiten.Image
}

func (imageSprite ImageSprite) Draw(image *ebiten.Image, center Point, rotation float64) {
	op := &ebiten.DrawImageOptions{}

	w := imageSprite.Image.Bounds().Max.X - imageSprite.Image.Bounds().Min.X
	h := imageSprite.Image.Bounds().Max.Y - imageSprite.Image.Bounds().Min.Y

	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(rotation)
	op.GeoM.Translate(center.x, center.y)

	image.DrawImage(imageSprite.Image, op)
}

type GameObject struct {
	id       int
	active   bool
	position Point
	rotation float64
}
