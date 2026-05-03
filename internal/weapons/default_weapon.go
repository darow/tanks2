package weapons

import (
	"math"
	"sync"
	"time"

	"myebiten/internal/models"
)

const (
	DEFAULT_GUN_BULLET_SPEED  = 1.15
	DEFAULT_GUN_BULLET_RADIUS = 4
	DEFAULT_GUN_BULLETS_COUNT = 3
	DEFAULT_GUN_BULLET_TTL    = 7 * time.Second
)

type DefaultWeapon struct {
	Clip         models.Pool[*models.Bullet]
	Cooldown     time.Duration
	BulletRadius float64
	BulletSpeed  float64
	mu           sync.Mutex
	lastShotAt   time.Time
}

func (dw *DefaultWeapon) Shoot(origin models.Vector2D, rotation float64) {
	dw.mu.Lock()
	if dw.Cooldown > 0 && !dw.lastShotAt.IsZero() && time.Since(dw.lastShotAt) < dw.Cooldown {
		dw.mu.Unlock()
		return
	}
	dw.lastShotAt = time.Now()
	dw.mu.Unlock()

	dw.spawnBullet(origin, rotation)
}

func (dw *DefaultWeapon) spawnBullet(origin models.Vector2D, rotation float64) {
	bullet := dw.Clip.Get()
	if bullet == nil {
		return
	}

	bullet.Position.X = origin.X
	bullet.Position.Y = origin.Y

	bullet.Rotation = rotation

	sin, cos := math.Sincos(rotation)
	bullet.Speed.X = cos * dw.bulletSpeed()
	bullet.Speed.Y = sin * dw.bulletSpeed()
	bullet.R = dw.bulletRadius()
	bullet.Sprite.R = dw.bulletRadius()

	bullet.SetActive(true)

	go func() {
		time.Sleep(DEFAULT_GUN_BULLET_TTL)
		bullet.SetActive(false)
	}()
}

func (dw *DefaultWeapon) bulletSpeed() float64 {
	if dw.BulletSpeed > 0 {
		return dw.BulletSpeed
	}

	return DEFAULT_GUN_BULLET_SPEED
}

func (dw *DefaultWeapon) bulletRadius() float64 {
	if dw.BulletRadius > 0 {
		return dw.BulletRadius
	}

	return DEFAULT_GUN_BULLET_RADIUS
}
