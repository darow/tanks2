package weapons

import (
	"time"

	"myebiten/internal/models"
)

func NewDefaultWeapon(clip models.Pool[*models.Bullet]) *DefaultWeapon {
	return &DefaultWeapon{
		Clip:     clip,
		Cooldown: time.Millisecond * 500,
	}
}
