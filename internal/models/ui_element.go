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
	font      font.Face
	Text      string
	TextColor color.RGBA
}

func (uiText UIText) Draw(drawingArea *DrawingArea) {
	textColor := uiText.TextColor
	if textColor.A == 0 {
		textColor = color.RGBA{0x00, 0x00, 0x00, 0xff}
	}

	text.Draw(drawingArea.BoardImage, uiText.Text, uiText.font, int(drawingArea.Offset.X), int(drawingArea.Offset.Y), textColor)
}

func (uiText *UIText) SetText(newText string) {
	uiText.Text = newText
}

func (uiText *UIText) SetColor(textColor color.RGBA) {
	uiText.TextColor = textColor
}

func CreateUIText(s string, font font.Face) UIText {
	return UIText{
		font:      font,
		Text:      s,
		TextColor: color.RGBA{0x00, 0x00, 0x00, 0xff},
	}
}

type UIPanel struct {
	UIElement
	sprite RectangleSprite
}

func (uiPanel UIPanel) Draw(drawingArea *DrawingArea) {
	uiPanel.sprite.Draw(0.0, 0.0, 0.0, drawingArea)
}
