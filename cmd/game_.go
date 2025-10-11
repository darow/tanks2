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

type Coordinates struct {
	i, j int
}

type MazeNode struct {
	up    bool
	down  bool
	right bool
	left  bool
}

func (mNode *MazeNode) addDirection(x, y int) {
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

func next(currentCoord, root Coordinates, N, M int) (Coordinates, bool) {
	if currentCoord == root {
		return Coordinates{}, false
	}

	dir := -1
	if root.i > currentCoord.i || (root.i == currentCoord.i && ((root.j > currentCoord.j && root.i%2 != 0) || (root.j < currentCoord.j && root.i%2 == 0))) {
		dir = 1
	}

	i, j := currentCoord.i, currentCoord.j

	if i%2 == 1 {
		if 1 <= j+dir && j+dir <= M {
			return Coordinates{i, j + dir}, true
		}

		return Coordinates{i + dir, j}, true
	}

	if 1 <= j-dir && j-dir <= M {
		return Coordinates{i, j - dir}, true
	}

	return Coordinates{i + dir, j}, true
}

func getRandomDirection(node Coordinates, prevDir, N, M int) int {
	distribution := [4]float32{0.25, 0.25, 0.25, 0.25}

	if prevDir != -1 {
		for i := range distribution {
			distribution[i] = 0.3
		}
		distribution[prevDir] = 0.1
	}

	p := rand.Float32()

	var s float32 = 0.0
	for index := range distribution {
		s += distribution[index]
		if p <= s {
			return index
		}
	}

	return -1
}

func getInitialMaze(N, M int) ([][]MazeNode, Coordinates) {
	mazeNodes := make([][]MazeNode, N+2)
	for i := range N + 2 {
		mazeNodes[i] = make([]MazeNode, M+2)
	}

	root := Coordinates{rand.Intn(N) + 1, rand.Intn(M) + 1}

	lastJ := M
	if N%2 == 0 {
		lastJ = 1
	}

	coords := [2]Coordinates{{1, 1}, {N, lastJ}}
	for _, coord := range coords {
		for {
			nextCoord, ok := next(coord, root, N, M)
			// fmt.Printf("%d %d\n", nextCoord.i, nextCoord.j)
			if !ok {
				break
			}

			mazeNodes[coord.i][coord.j].addDirection(nextCoord.j-coord.j, nextCoord.i-coord.i)
			coord = nextCoord

			// once you joined existing branch you can quit
			if node := mazeNodes[coord.i][coord.j]; node.up || node.down || node.right || node.left {
				break
			}
		}
	}

	return mazeNodes, root
}

func addConnections(mazeNodes [][]MazeNode) [][]MazeNode {
	// count := min(len(mazeNodes), len(mazeNodes[0])) - 2

	return mazeNodes
}

func createMaze(N, M int) [][]MazeNode {
	mazeNodes, root := getInitialMaze(N, M)

	dirIndex := -1
	count := 0
	for count < N*M {
		dirIndex = getRandomDirection(root, dirIndex, N, M)

		switch dirIndex {
		case 0:
			if root.j == M {
				continue
			}
			mazeNodes[root.i][root.j].addDirection(1, 0)
			root.j++
			count++
		case 1:
			if root.i == N {
				continue
			}
			mazeNodes[root.i][root.j].addDirection(0, 1)
			root.i++
			count++
		case 2:
			if root.j == 1 {
				continue
			}
			mazeNodes[root.i][root.j].addDirection(-1, 0)
			root.j--
			count++
		case 3:
			if root.i == 1 {
				continue
			}
			mazeNodes[root.i][root.j].addDirection(0, -1)
			root.i--
			count++
		}

		node := &mazeNodes[root.i][root.j]
		if node.up {
			node.up = false
			continue
		}
		if node.down {
			node.down = false
			continue
		}
		if node.right {
			node.right = false
			continue
		}
		if node.left {
			node.left = false
			continue
		}
	}

	mazeNodes = addConnections(mazeNodes)

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
