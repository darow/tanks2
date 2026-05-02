package weapons

import (
	"time"

	"myebiten/internal/models"
)

type RocketWeapon struct {
	DefaultWeapon
}

func NewRocketWeapon(clip models.Pool[*models.Bullet]) *RocketWeapon {
	return &RocketWeapon{
		DefaultWeapon: DefaultWeapon{
			Clip:         clip,
			Cooldown:     900 * time.Millisecond,
			BulletRadius: 12,
			BulletSpeed:  1.35,
		},
	}
}
