package game

import (
	"bytes"
	"image"
	"log"
	"math/rand"

	"myebiten/internal/models"
	"myebiten/internal/models/item"
	images "myebiten/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

const itemIconSize = 34

var itemSpriteBytes = [][]byte{
	images.ExplosionPng,
	images.MinigunPng,
	images.RocketPng,
}

type itemSprite struct {
	itemType item.ItemType
	sprite   *models.ImageSprite
}

var itemSprites []itemSprite

func loadItemSprites() []itemSprite {
	sprites := make([]itemSprite, 0, len(itemSpriteBytes))

	for idx, raw := range itemSpriteBytes {
		img, _, err := image.Decode(bytes.NewReader(raw))
		if err != nil {
			log.Println(err)
			continue
		}

		resized := resize.Resize(itemIconSize, 0, img, resize.Lanczos3)
		sprite := itemSprite{
			itemType: item.ItemType(idx),
			sprite:   &models.ImageSprite{Image: ebiten.NewImageFromImage(resized)},
		}
		sprites = append(sprites, sprite)
	}

	return sprites
}

func getItemSprite(itemType item.ItemType) *models.ImageSprite {
	if len(itemSprites) == 0 {
		itemSprites = loadItemSprites()
	}

	for _, itemSprite := range itemSprites {
		if itemSprite.itemType == itemType {
			return itemSprite.sprite
		}
	}

	return nil
}

func getRandomItemSprite() (item.ItemType, *models.ImageSprite) {
	if len(itemSprites) == 0 {
		itemSprites = loadItemSprites()
	}

	selected := itemSprites[rand.Intn(len(itemSprites))]

	//TODO: REMOVE THIS LINES TO SPAWN OTHER ITEM TYPES
	selected.itemType = item.TypeMinigun
	selected.sprite = getItemSprite(item.TypeMinigun)

	return selected.itemType, selected.sprite
}
