package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

var (
	wd                WallsDTO
	sendingMapUpdates bool
)

func (g *Game) CreateMap() {
	g.Things = Things{
		Bullets: make(map[int]Bullet),
		walls:   make(map[Wall]struct{}),
	}

	for x := range g.boardSizeX {
		w1 := Wall{
			X:          uint16(x),
			Y:          uint16(0),
			Horizontal: true,
		}
		w2 := Wall{
			X:          uint16(x),
			Y:          uint16(g.boardSizeY),
			Horizontal: true,
		}
		g.Things.walls[w1] = struct{}{}
		g.Things.walls[w2] = struct{}{}
	}

	for y := range g.boardSizeY {
		w1 := Wall{
			X:          uint16(0),
			Y:          uint16(y),
			Horizontal: false,
		}
		w2 := Wall{
			X:          uint16(g.boardSizeX),
			Y:          uint16(y),
			Horizontal: false,
		}
		g.Things.walls[w1] = struct{}{}
		g.Things.walls[w2] = struct{}{}
	}

	for y := 1; y < g.boardSizeY; y++ {
		for x := 1; x < g.boardSizeX; x++ {
			n := rand.Int()
			//generate horizontal
			if x < g.boardSizeX-1 && n%100 < 25 {
				w := Wall{
					X:          uint16(x),
					Y:          uint16(y),
					Horizontal: true,
				}
				g.Things.walls[w] = struct{}{}
			}

			//generate vertical
			if y < g.boardSizeY-1 && n%100 < 45 {
				w := Wall{
					X:          uint16(x),
					Y:          uint16(y),
					Horizontal: false,
				}
				g.Things.walls[w] = struct{}{}
			}
		}
	}

	for _, char := range g.CharactersStash {
		g.Characters[char.id] = char
	}
	g.CharactersStash = nil

	spawnPlaces := []Point{
		Point{x: WALL_HEIGHT * 0.5, y: WALL_HEIGHT * 0.5},
		Point{x: WALL_HEIGHT*(float64(g.boardSizeX)-1) + WALL_HEIGHT*0.5, y: WALL_HEIGHT*(float64(g.boardSizeY)-1) + WALL_HEIGHT*0.5},
	}

	i := 0
	for _, char := range g.Characters {
		if char == nil {
			continue
		}

		char.X = spawnPlaces[i].x
		char.Y = spawnPlaces[i].y
		i++
	}
}

func (g *Game) SendMapToClient() {
	wd = WallsDTO{}

	for key := range g.Things.walls {
		wd.Walls = append(wd.Walls, key)
	}

	msg, err := json.Marshal(wd)
	if err != nil {
		log.Fatal(err)
	}

	err = g.server.WriteMapMessage(msg)

	if err != nil {
		log.Fatal(err)
	}

	//if !sendingMapUpdates {
	//	go g.sendMapUpdates()
	//	sendingMapUpdates = true
	//}
}

func (g *Game) sendMapUpdates() {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		msg, err := json.Marshal(wd)
		if err != nil {
			log.Fatal(err)
		}

		err = g.server.WriteMapMessage(msg)

		if err != nil {
			log.Fatal(err)
		}
	}
}
