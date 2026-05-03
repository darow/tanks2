package game

import (
	"myebiten/internal/models"
	"myebiten/internal/models/character"
	"myebiten/internal/models/item"
	"myebiten/internal/weapons"
)

func (mainScene *MainScene) applyItemEffect(itemToApply *item.Item, char *character.Character, charIndex int) {
	switch itemToApply.Type {
	case item.TypeExplosion:
		mainScene.applyExplosion(char, charIndex)
	case item.TypeMinigun:
		mainScene.applyMinigun(char)
	case item.TypeRocket:
		mainScene.applyRocket(char, charIndex)
	}
}

func (mainScene *MainScene) applyExplosion(char *character.Character, charIndex int) {
	clip := mainScene.weaponClipFor(charIndex)
	char.SetWeapon(weapons.NewExplosionWeapon(clip))
}

func (mainScene *MainScene) applyMinigun(char *character.Character) {
	start := weapons.DEFAULT_GUN_BULLETS_COUNT*mainScene.PlayersCount + char.ID*weapons.MINIGUN_BULLETS_COUNT
	end := weapons.DEFAULT_GUN_BULLETS_COUNT*mainScene.PlayersCount + (char.ID+1)*weapons.MINIGUN_BULLETS_COUNT
	clip := models.CreatePool(mainScene.Bullets[start:end])

	char.SetWeapon(weapons.NewMinigunWeapon(clip))
}

func (mainScene *MainScene) applyRocket(char *character.Character, charIndex int) {
	clip := mainScene.weaponClipFor(charIndex)
	char.SetWeapon(weapons.NewRocketWeapon(clip))
}

func (mainScene *MainScene) weaponClipFor(charIndex int) models.Pool[*models.Bullet] {
	start := charIndex * weapons.DEFAULT_GUN_BULLETS_COUNT
	end := (charIndex + 1) * weapons.DEFAULT_GUN_BULLETS_COUNT
	return models.CreatePool(mainScene.Bullets[start:end])
}
