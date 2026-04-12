package game

import (
	"bytes"
	"image"
	"log"
	"math/rand"

	"myebiten/internal/models"
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

var itemSprites []*models.ImageSprite

func getRandomItemSprite() *models.ImageSprite {
	if len(itemSprites) == 0 {
		itemSprites = loadItemSprites()
	}

	if len(itemSprites) == 0 {
		return nil
	}

	return itemSprites[rand.Intn(len(itemSprites))]
}

func loadItemSprites() []*models.ImageSprite {
	sprites := make([]*models.ImageSprite, 0, len(itemSpriteBytes))

	for _, raw := range itemSpriteBytes {
		img, _, err := image.Decode(bytes.NewReader(raw))
		if err != nil {
			log.Println(err)
			continue
		}

		resized := resize.Resize(itemIconSize, 0, img, resize.Lanczos3)
		sprite := &models.ImageSprite{Image: ebiten.NewImageFromImage(resized)}
		sprites = append(sprites, sprite)
	}

	return sprites
}
