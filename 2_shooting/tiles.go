package main

import (
	"log"
	"math"
)

const (
	BULLET_RADIUS = 4
)

var TILE_ID_SEQUENCE = uint16(0)

type Tiles struct {
	bullets map[uint16]Bullet
}

type Bullet struct {
	id   uint16
	x, y float32
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
