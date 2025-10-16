package main

import "github.com/hajimehoshi/ebiten/v2"

type DrawingSettings struct {
	Offset Vector2D
	Scale  float64
}

type DrawingArea struct {
	DrawingSettings

	boardImage *ebiten.Image

	Height float64
	Width  float64

	parent *DrawingArea
}

func (drawingArea *DrawingArea) NewArea(height, width float64, settings DrawingSettings) (newArea *DrawingArea) {
	newArea = &DrawingArea{
		Height: height,
		Width:  width,

		boardImage: drawingArea.boardImage,

		DrawingSettings: DrawingSettings{
			Offset: Vector2D{
				x: drawingArea.Offset.x + settings.Offset.x,
				y: drawingArea.Offset.y + settings.Offset.y,
			},
			Scale: drawingArea.Scale * settings.Scale,
		},

		parent: drawingArea,
	}

	return
}
