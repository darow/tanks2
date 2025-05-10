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

//type WallCoords struct {
//	x, y uint16
//}
