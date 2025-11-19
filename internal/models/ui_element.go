package models

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type UIElement struct {
	Active bool
}

func (uiElem UIElement) IsActive() bool {
	return uiElem.Active
}

type UIText struct {
	UIElement
	font font.Face
	text string
}

func (uiText UIText) Draw(drawingArea *DrawingArea) {
	text.Draw(drawingArea.BoardImage, uiText.text, uiText.font, int(drawingArea.Offset.X), int(drawingArea.Offset.Y), color.Black)
}

func (uiText *UIText) SetText(newText string) {
	uiText.text = newText
}
