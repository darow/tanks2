package models

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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

func (gameObject *GameObject) Move() {
	gameObject.Position.X += gameObject.Speed.X
	gameObject.Position.Y += gameObject.Speed.Y
}

func (gameObject *GameObject) MoveBack() {
	gameObject.Position.X -= gameObject.Speed.X
	gameObject.Position.Y -= gameObject.Speed.Y
}

type DrawingArea struct {
	DrawingSettings

	BoardImage *ebiten.Image

	Height float64
	Width  float64

	Children []*DrawingArea
}

type DrawingSettings struct {
	Offset Vector2D
	Scale  float64
}

func (d *DrawingArea) NewArea(height, width float64, settings DrawingSettings) (newArea *DrawingArea) {
	newArea = &DrawingArea{
		Height: height,
		Width:  width,

		BoardImage: d.BoardImage,

		DrawingSettings: DrawingSettings{
			Offset: Vector2D{
				X: d.Offset.X + settings.Offset.X,
				Y: d.Offset.Y + settings.Offset.Y,
			},
			Scale: d.Scale * settings.Scale,
		},
	}

	d.Children = append(d.Children, newArea)
	return
}
