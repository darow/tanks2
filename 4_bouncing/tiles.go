package main

import (
	"log"
	"math"
)

const (
	BULLET_RADIUS            = 4
	CHARACTER_ROTATION_SPEED = 0.05
	CHARACTER_SPEED          = 3
	BULLET_SPEED             = 4
	WALL_HEIGHT              = 200
	WALL_WIDTH               = 10
)

var TILE_ID_SEQUENCE = uint16(0)

type Tiles struct {
	bullets map[uint16]Bullet
	walls   map[Wall]struct{}
}

type Bullet struct {
	id       uint16
	x, y     float32
	rotation float64
}

func (t *Tiles) getNextID() uint16 {
	TILE_ID_SEQUENCE++

	if len(t.bullets) >= math.MaxUint16 {
		log.Fatal("can not get next id: all possible values are used. Increase id type to uint32")
	}

	if _, ok := t.bullets[TILE_ID_SEQUENCE]; ok {
		return t.getNextID()
	}

	return TILE_ID_SEQUENCE
}

type Character struct {
	x, y     float32
	rotation float64
}

type Wall struct {
	x, y       uint16
	horizontal bool
}

func (t *Tiles) ProcessBulletToWallCollision(b Bullet, dx, dy float32) Bullet {
	isCollision, isHorizontal := t.DetectBulletToWallCollision(b, dx, dy)

	if isCollision {
		if isHorizontal {
			b.rotation = -b.rotation

			return b
		}

		if b.rotation < math.Pi {
			b.rotation = math.Pi - b.rotation
		} else {
			b.rotation = b.rotation - math.Pi
		}
	}

	return b
}

func (t *Tiles) DetectBulletToWallCollision(b Bullet, dx, dy float32) (isCollision, isHorizontal bool) {
	if int(b.x)%WALL_HEIGHT <= WALL_WIDTH {
		wallToCollide := Wall{
			x:          uint16(math.Floor(float64(b.x / WALL_HEIGHT))),
			y:          uint16(math.Floor(float64(b.y / WALL_HEIGHT))),
			horizontal: false,
		}

		if _, ok := t.walls[wallToCollide]; ok {
			isCollision, isHorizontal = true, false
			//end of wall
			if math.Abs(float64(int(b.y+float32(math.Abs(float64(dy))))%WALL_HEIGHT)) <= WALL_WIDTH {
				isHorizontal = true
			}

			return
		}
	}

	if int(b.y)%WALL_HEIGHT <= WALL_WIDTH {
		wallToCollide := Wall{
			x:          uint16(math.Floor(float64(b.x / WALL_HEIGHT))),
			y:          uint16(math.Floor(float64(b.y / WALL_HEIGHT))),
			horizontal: true,
		}

		if _, ok := t.walls[wallToCollide]; ok {
			isCollision, isHorizontal = true, true

			//end of wall
			if math.Abs(float64(int(b.x+float32(math.Abs(float64(dx))))%WALL_HEIGHT)) <= WALL_WIDTH {
				isHorizontal = false
			}

			return
		}
	}

	return
}
