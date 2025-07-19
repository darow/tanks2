package main

import (
	"math/rand"
)

func (g *Game) CreateMap() {
	g.things = Things{
		bullets: make(map[int]Bullet),
		walls:   make(map[Wall]struct{}),
	}

	for x := range g.boardSizeX {
		w1 := Wall{
			x:          uint16(x),
			y:          uint16(0),
			horizontal: true,
		}
		w2 := Wall{
			x:          uint16(x),
			y:          uint16(g.boardSizeY),
			horizontal: true,
		}
		g.things.walls[w1] = struct{}{}
		g.things.walls[w2] = struct{}{}
	}

	for y := range g.boardSizeY {
		w1 := Wall{
			x:          uint16(0),
			y:          uint16(y),
			horizontal: false,
		}
		w2 := Wall{
			x:          uint16(g.boardSizeX),
			y:          uint16(y),
			horizontal: false,
		}
		g.things.walls[w1] = struct{}{}
		g.things.walls[w2] = struct{}{}
	}

	for y := 1; y < g.boardSizeY; y++ {
		for x := 1; x < g.boardSizeX; x++ {
			n := rand.Int()
			//generate horizontal
			if x < g.boardSizeX-1 && n%100 < 25 {
				w := Wall{
					x:          uint16(x),
					y:          uint16(y),
					horizontal: true,
				}
				g.things.walls[w] = struct{}{}
			}

			//generate vertical
			if y < g.boardSizeY-1 && n%100 < 45 {
				w := Wall{
					x:          uint16(x),
					y:          uint16(y),
					horizontal: false,
				}
				g.things.walls[w] = struct{}{}
			}
		}
	}

	for _, char := range g.charactersStash {
		g.characters[char.id] = char
	}
	g.charactersStash = nil

	spawnPlaces := []Point{
		Point{x: WALL_HEIGHT * 0.5, y: WALL_HEIGHT * 0.5},
		Point{x: WALL_HEIGHT*(float64(g.boardSizeX)-1) + WALL_HEIGHT*0.5, y: WALL_HEIGHT*(float64(g.boardSizeY)-1) + WALL_HEIGHT*0.5},
	}

	i := 0
	for _, char := range g.characters {
		if char == nil {
			continue
		}

		char.x = spawnPlaces[i].x
		char.y = spawnPlaces[i].y
		i++
	}
}
