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

func (uiElem *UIElement) SetActive(b bool) {
	uiElem.Active = b
}

type UIText struct {
	UIElement
	font font.Face
	Text string
}

func (uiText UIText) Draw(drawingArea *DrawingArea) {
	text.Draw(drawingArea.BoardImage, uiText.Text, uiText.font, int(drawingArea.Offset.X), int(drawingArea.Offset.Y), color.Black)
}

func (uiText *UIText) SetText(newText string) {
	uiText.Text = newText
}

func CreateUIText(s string, font font.Face) UIText {
	return UIText{
		font: font,
		Text: s,
	}
}

type UIPanel struct {
	UIElement
	sprite RectangleSprite
}

func (uiPanel UIPanel) Draw(drawingArea *DrawingArea) {
	uiPanel.sprite.Draw(0.0, 0.0, 0.0, drawingArea)
}
