package models

import "image/color"

const ITEM_SIZE = 56

type Item struct {
	GameObject
	BaseSprite RectangleSprite
	IconSprite ImageSprite
}

func (item *Item) Draw(drawingArea *DrawingArea) {
	if drawingArea == nil {
		return
	}

	drawFilledRect(
		drawingArea,
		item.Position.X,
		item.Position.Y,
		item.BaseSprite.W,
		item.BaseSprite.H,
		color.RGBA{0x90, 0x90, 0x90, 0xff},
	)

	item.IconSprite.Draw(item.Position.X, item.Position.Y, 0.0, drawingArea)
}

func CreateItem(position Vector2D, iconImage *ImageSprite) *Item {
	item := &Item{
		GameObject: GameObject{
			Position: position,
			Active:   true,
		},
		BaseSprite: RectangleSprite{
			W: ITEM_SIZE,
			H: ITEM_SIZE,
		},
	}

	if iconImage != nil {
		item.IconSprite = *iconImage
	}

	return item
}
