package models

import (
	"math"
	"time"
)

const (
	BULLET_SPEED  = 6
	BULLETS_COUNT = 5
)

type Weapon interface {
	Shoot(origin Vector2D, rotation float64)
	Discharge()
}

type DefaultWeapon struct {
	Clip     []*Bullet
	Cooldown int64
}

func (dw *DefaultWeapon) Shoot(origin Vector2D, rotation float64) {
	for _, bullet := range dw.Clip {
		if !bullet.IsActive() {
			bullet.Position.X = origin.X
			bullet.Position.Y = origin.Y

			bullet.Rotation = rotation

			sin, cos := math.Sincos(rotation)
			bullet.Speed.X = cos * BULLET_SPEED
			bullet.Speed.Y = sin * BULLET_SPEED

			bullet.SetActive(true)

			go func() {
				time.Sleep(time.Duration(dw.Cooldown) * time.Second)
				bullet.SetActive(false)
			}()

			break
		}
	}

	dw.Discharge()
}

func (dw *DefaultWeapon) Discharge() {

}
