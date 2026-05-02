package weapons

import (
	"time"

	"myebiten/internal/models"
)

type ExplosionWeapon struct {
	DefaultWeapon
}

func NewExplosionWeapon(clip models.Pool[*models.Bullet]) *ExplosionWeapon {
	return &ExplosionWeapon{
		DefaultWeapon: DefaultWeapon{
			Clip:         clip,
			Cooldown:     250 * time.Millisecond,
			BulletRadius: 18,
			BulletSpeed:  0.7,
		},
	}
}
