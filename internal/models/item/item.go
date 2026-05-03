package item

import (
	"myebiten/internal/models"
	"myebiten/internal/models/character"
)

const ITEM_SIZE = 56

type ItemType int

const (
	TypeExplosion ItemType = iota
	TypeMinigun
	TypeRocket
)

type Item struct {
	models.GameObject
	Type       ItemType
	IconSprite models.ImageSprite `json:"-"`
}

func (item *Item) Draw(drawingArea *models.DrawingArea) {
	if drawingArea == nil {
		return
	}

	item.IconSprite.Draw(item.Position.X, item.Position.Y, 0.0, drawingArea)
}

func (item *Item) DetectCharacterCollision(char *character.Character) bool {
	if item == nil || char == nil || !item.IsActive() || !char.IsActive() {
		return false
	}

	pickupRadius := (float64(ITEM_SIZE) + float64(character.CHARACTER_WIDTH)) / 2
	return models.SquareDistance(item.Position, char.Position) <= pickupRadius*pickupRadius
}

func CreateItem(itemType ItemType, position models.Vector2D, iconImage *models.ImageSprite) *Item {
	item := &Item{
		GameObject: models.GameObject{
			Position: position,
			Active:   true,
		},
		Type: itemType,
	}

	if iconImage != nil {
		item.IconSprite = *iconImage
	}

	return item
}
