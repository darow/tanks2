package weapons

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"myebiten/internal/models"
)

type MinigunWeapon struct {
	DefaultWeapon
	mu            sync.Mutex
	isShooting    bool
	shootHeldTill time.Time
	origin        models.Vector2D
	rotation      float64
}

const (
	MINIGUN_WARMUP            = 500 * time.Millisecond
	MINIGUN_BULLETS_COUNT     = 30
	MINIGUN_DISPERSION_DEGREE = 10.0
	MINIGUN_COOLDOWN          = 100 * time.Millisecond
)

func NewMinigunWeapon(clip models.Pool[*models.Bullet]) *MinigunWeapon {
	return &MinigunWeapon{
		DefaultWeapon: DefaultWeapon{
			Clip:         clip,
			Cooldown:     MINIGUN_COOLDOWN,
			BulletRadius: 2,
			BulletSpeed:  DEFAULT_GUN_BULLET_SPEED * 1.1,
		},
	}
}

func (mw *MinigunWeapon) Shoot(origin models.Vector2D, rotation float64) {
	mw.mu.Lock()
	mw.shootHeldTill = time.Now().Add(MINIGUN_COOLDOWN)
	mw.origin = origin
	mw.rotation = rotation
	if mw.isShooting {
		mw.mu.Unlock()
		return
	}
	mw.isShooting = true
	mw.mu.Unlock()

	go mw.fireBurst(origin, rotation)
}

func (mw *MinigunWeapon) IsShooting() bool {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	return mw.isShooting
}

func (mw *MinigunWeapon) fireBurst(origin models.Vector2D, rotation float64) {
	timer := time.NewTimer(MINIGUN_WARMUP)
	defer timer.Stop()

	<-timer.C

	for i := 0; i < MINIGUN_BULLETS_COUNT; i++ {
		mw.mu.Lock()
		held := time.Now().Before(mw.shootHeldTill)
		origin := mw.origin
		rotation := mw.rotation
		if !held {
			mw.isShooting = false
			mw.mu.Unlock()
			return
		}
		mw.mu.Unlock()

		dispersion := (rand.Float64()*2 - 1) * MINIGUN_DISPERSION_DEGREE * math.Pi / 180
		mw.spawnBullet(origin, rotation+dispersion)

		time.Sleep(mw.Cooldown)
	}

	mw.mu.Lock()
	mw.isShooting = false
	mw.mu.Unlock()
}
