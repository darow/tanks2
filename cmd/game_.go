package main

import (
	"encoding/json"
	"log"
	"time"
)

var (
	wd                WallsDTO
	sendingMapUpdates bool
)

type MazeNode struct {
	up    bool
	down  bool
	right bool
	left  bool
}

func (mNode *MazeNode) addDirection(y, x int) {
	if y == 0 {
		if x == 1 {
			mNode.right = true
		} else {
			mNode.left = true
		}
	} else if y == -1 {
		mNode.down = true
	} else {
		mNode.up = true
	}
}

func (g *Game) CreateMap() {
	g.Things = Things{
		Bullets: make(map[int]Bullet),
		walls:   make(map[Wall]struct{}),
	}

	mazeNodes := createMaze(g.boardSizeY, g.boardSizeX)
	buildMaze(mazeNodes, g.Things.walls)

	for _, char := range g.CharactersStash {
		g.Characters[char.id] = char
	}
	g.CharactersStash = nil

	spawnPlaces := []Point{
		{x: WALL_HEIGHT * 0.5, y: WALL_HEIGHT * 0.5},
		{x: WALL_HEIGHT*(float64(g.boardSizeX)-1) + WALL_HEIGHT*0.5, y: WALL_HEIGHT*(float64(g.boardSizeY)-1) + WALL_HEIGHT*0.5},
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

func next(i, j, N, M int) (int, int, bool) {
	if (i == N) && (((N%2 == 0) && (j == 1)) || ((N%2 == 1) && (j == M))) {
		return i, j, false
	}

	if i%2 == 1 {
		if j < M {
			return i, j + 1, true
		}

		return i + 1, j, true
	}

	if j > 1 {
		return i, j - 1, true
	}

	return i + 1, j, true
}

func getInitialMaze(N, M int) [][]MazeNode {
	mazeNodes := make([][]MazeNode, N+2)

	for i := range N + 2 {
		mazeNodes[i] = make([]MazeNode, M+2)
	}

	i, j := 1, 1
	for {
		i1, j1, ok := next(i, j, N, M)
		// fmt.Printf("%d %d\n", i1, j1)
		if !ok {
			break
		}

		mazeNodes[i][j].addDirection(i1-i, j1-j)
		mazeNodes[i1][j1].addDirection(i-i1, j-j1)
		i, j = i1, j1
	}

	return mazeNodes
}

func createMaze(N, M int) [][]MazeNode {
	mazeNodes := getInitialMaze(N, M)

	return mazeNodes
}

func buildMaze(mazeNodes [][]MazeNode, walls map[Wall]struct{}) {
	for i := 1; i < len(mazeNodes); i++ {
		for j := 1; j < len(mazeNodes[0]); j++ {
			currentNode := mazeNodes[i][j]
			leftNode := mazeNodes[i][j-1]
			downNode := mazeNodes[i-1][j]

			horizontalWall := !(currentNode.down || downNode.up) && (j != len(mazeNodes[0])-1)
			verticalWall := !(currentNode.left || leftNode.right) && (i != len(mazeNodes)-1)

			if horizontalWall {
				w := Wall{
					X:          uint16(j - 1),
					Y:          uint16(i - 1),
					Horizontal: true,
				}
				walls[w] = struct{}{}
			}

			if verticalWall {
				w := Wall{
					X:          uint16(j - 1),
					Y:          uint16(i - 1),
					Horizontal: false,
				}
				walls[w] = struct{}{}
			}
		}
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
