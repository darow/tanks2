package models

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var COLOR_BACKGROUND = color.RGBA{0xca, 0xca, 0xff, 0xff}

type Drawable interface {
	IsActive() bool
	Draw(*DrawingArea)
}

type DrawingSettings struct {
	Offset Vector2D
	Scale  float64
}

type DrawingArea struct {
	DrawingSettings

	BoardImage *ebiten.Image

	Height float64
	Width  float64

	Children []*DrawingArea
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

type Scene struct {
	Objects  []Drawable
	rootArea *DrawingArea
	areaIDs  map[Drawable]string
	areas    map[string]*DrawingArea
}

func (scene *Scene) Draw() *ebiten.Image {
	boardImage := scene.rootArea.BoardImage
	boardImage.Clear()
	boardImage.Fill(COLOR_BACKGROUND)

	for _, object := range scene.Objects {
		if object.IsActive() {
			areaID := scene.areaIDs[object]
			area := scene.areas[areaID]
			object.Draw(area)
		}
	}

	return boardImage
}

func (scene *Scene) AddObject(object Drawable, areaID string) {
	scene.Objects = append(scene.Objects, object)
	scene.areaIDs[object] = areaID
}

func (scene *Scene) AddDrawingArea(areaID string, drawingArea *DrawingArea) {
	scene.areas[areaID] = drawingArea
}

func (scene *Scene) GetRootArea() *DrawingArea {
	return scene.rootArea
}

func (scene *Scene) GetArea(areaID string) *DrawingArea {
	return scene.areas[areaID]
}

func CreateScene(image *ebiten.Image, height, width float64) *Scene {
	area := &DrawingArea{
		BoardImage: image,

		DrawingSettings: DrawingSettings{
			Scale: 1.0,
		},

		Height: height,
		Width:  width,
	}

	return &Scene{
		rootArea: area,
		areas:    map[string]*DrawingArea{"root_area": area},
		areaIDs:  map[Drawable]string{},
	}
}
