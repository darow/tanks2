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

type SceneUI struct {
	Objects  []Drawable
	rootArea *DrawingArea
	AreaIDs  map[Drawable]string
	Areas    map[string]*DrawingArea
}

func (sceneUI *SceneUI) Draw() *ebiten.Image {
	boardImage := sceneUI.rootArea.BoardImage
	boardImage.Clear()
	boardImage.Fill(COLOR_BACKGROUND)

	for _, object := range sceneUI.Objects {
		if object.IsActive() {
			areaID := sceneUI.AreaIDs[object]
			area := sceneUI.Areas[areaID]
			object.Draw(area)
		}
	}

	return boardImage
}

func (sceneUI *SceneUI) AddObject(object Drawable, areaID string) {
	sceneUI.Objects = append(sceneUI.Objects, object)
	sceneUI.AreaIDs[object] = areaID
}

func (sceneUI *SceneUI) AddDrawingArea(areaID string, drawingArea *DrawingArea) {
	sceneUI.Areas[areaID] = drawingArea
}

func (sceneUI *SceneUI) GetRootArea() *DrawingArea {
	return sceneUI.rootArea
}

func (sceneUI *SceneUI) GetArea(areaID string) *DrawingArea {
	return sceneUI.Areas[areaID]
}

type Scene interface {
	Update() error
	Draw() *ebiten.Image
}

func CreateSceneUI(image *ebiten.Image, height, width float64) SceneUI {
	area := &DrawingArea{
		BoardImage: image,

		DrawingSettings: DrawingSettings{
			Scale: 1.0,
		},

		Height: height,
		Width:  width,
	}

	return SceneUI{
		rootArea: area,
		Areas:    map[string]*DrawingArea{"root_area": area},
		AreaIDs:  map[Drawable]string{},
	}
}
