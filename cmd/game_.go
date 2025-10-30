package main

import (
	"math"
	"math/rand"
)

type Coordinates struct {
	i, j int
}

type MazeNode struct {
	up    bool
	down  bool
	right bool
	left  bool

	topWall, bottomWall, rightWall, leftWall *Wall
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

func getSceneCoordinates(i, j int) Vector2D {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	return Vector2D{float64(j-1)*(wh-ww) + wh/2, float64(i-1)*(wh-ww) + wh/2}
}

func (g *Game) SetupLevel() (int, int, []Wall) {
	h := rand.Intn(MAX_BOARD_HEIGHT-MIN_BOARD_HEIGHT) + MIN_BOARD_HEIGHT
	w := rand.Intn(MAX_BOARD_WIDTH-MIN_BOARD_WIDTH) + MIN_BOARD_WIDTH

	walls := g.CreateMaze(h, w)
	g.SetDrawingSettings(h, w)
	g.SetCharacters(h, w)

	return h, w, walls
}

func (g *Game) SetCharacters(h, w int) {
	spawnPlaces := []Vector2D{}
	for range g.Characters {
		i := rand.Intn(h) + 1
		j := rand.Intn(w) + 1
		spawnPlace := getSceneCoordinates(i, j)
		spawnPlaces = append(spawnPlaces, spawnPlace)
	}

	i := 0
	for _, char := range g.Characters {
		if !char.Active {
			continue
		}

		char.Position.X = spawnPlaces[i].X
		char.Position.Y = spawnPlaces[i].Y

		char.Rotation = 0.0

		char.Speed.X = 0
		char.Speed.Y = 0

		i++
	}
}

func (g *Game) CreateMaze(h, w int) []Wall {
	g.Walls = make([]Wall, 0)

	g.Maze = createMaze(h, w)
	g.Walls = buildMaze(g.Maze, g.Walls)

	return g.Walls
}

func (g *Game) SetDrawingSettings(h, w int) {
	areaHeight := g.mainArea.Height
	areaWidth := g.mainArea.Width

	mazeHeight := float64(h*(WALL_HEIGHT-WALL_WIDTH) + WALL_WIDTH)
	mazeWidth := float64(w*(WALL_HEIGHT-WALL_WIDTH) + WALL_WIDTH)

	scalingFactor := min(areaHeight/mazeHeight, areaWidth/mazeWidth)

	mazeHeight *= scalingFactor
	mazeWidth *= scalingFactor

	newDrawingSettings := DrawingSettings{
		Offset: Vector2D{(areaWidth - mazeWidth) / 2, (areaHeight - mazeHeight) / 2},
		Scale:  scalingFactor,
	}
	newMainArea := g.mainArea.NewArea(mazeHeight, mazeWidth, newDrawingSettings)
	g.mainArea = newMainArea
}

func getRandomDirection(prevDir int) int {
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

	randomInt := rand.Intn(len(Generators))

	sources := Generators[randomInt].sources
	next := Generators[randomInt].next

	coords := sources(N, M)
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
	count := min(len(mazeNodes), len(mazeNodes[0])) - 2
	total := (len(mazeNodes)-2)*(len(mazeNodes[0])-2) - len(mazeNodes) - len(mazeNodes[0]) + 4
	p := float64(count) / float64(total)

	for i := 1; i <= len(mazeNodes)-2; i++ {
		for j := 1; j <= len(mazeNodes[0])-2; j++ {
			randomFloat := rand.Float64()
			if i != len(mazeNodes)-2 && !mazeNodes[i][j].up && randomFloat <= p {
				mazeNodes[i][j].up = true
			}

			randomFloat = rand.Float64()
			if j != len(mazeNodes[0])-2 && !mazeNodes[i][j].right && randomFloat <= p {
				mazeNodes[i][j].right = true
			}
		}
	}

	return mazeNodes
}

func createMaze(N, M int) [][]MazeNode {
	mazeNodes, root := getInitialMaze(N, M)

	dirIndex := -1
	count := 0
	for count < N*M {
		dirIndex = getRandomDirection(dirIndex)

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

func buildMaze(mazeNodes [][]MazeNode, walls []Wall) []Wall {
	for i := 1; i < len(mazeNodes); i++ {
		for j := 1; j < len(mazeNodes[0]); j++ {
			currentNode := &mazeNodes[i][j]
			leftNode := &mazeNodes[i][j-1]
			downNode := &mazeNodes[i-1][j]

			horizontalWall := !(currentNode.down || downNode.up) && (j != len(mazeNodes[0])-1)
			verticalWall := !(currentNode.left || leftNode.right) && (i != len(mazeNodes)-1)

			wh := float64(WALL_HEIGHT)
			ww := float64(WALL_WIDTH)

			nodeCenter := getSceneCoordinates(i, j)

			if horizontalWall {
				w := Wall{
					GameObject: GameObject{
						Active:   true,
						Position: Vector2D{nodeCenter.X, nodeCenter.Y - (wh-ww)/2},
						Rotation: 0.0,
					},
					Hitbox: RectangleHitbox{WALL_WIDTH, WALL_HEIGHT},
					Sprite: RectangleSprite{WALL_WIDTH, WALL_HEIGHT},
				}
				currentNode.bottomWall = &w
				downNode.topWall = &w
				walls = append(walls, w)
			}

			if verticalWall {
				w := Wall{
					GameObject: GameObject{
						Active:   true,
						Position: Vector2D{nodeCenter.X - (wh-ww)/2, nodeCenter.Y},
						Rotation: math.Pi / 2,
					},
					Hitbox: RectangleHitbox{WALL_WIDTH, WALL_HEIGHT},
					Sprite: RectangleSprite{WALL_WIDTH, WALL_HEIGHT},
				}
				walls = append(walls, w)
				currentNode.leftWall = &w
				leftNode.rightWall = &w
				walls = append(walls, w)
			}
		}
	}

	return walls
}

func (g *Game) Reset() {
	for _, bullet := range g.Bullets {
		bullet.Active = false
	}

	for _, char := range g.Characters {
		char.Active = true

		char.input.RotateLeft = false
		char.input.RotateRight = false
		char.input.MoveBackward = false
		char.input.MoveForward = false
		char.input.Shoot = false
	}

	if g.mainArea.parent != nil {
		g.mainArea = g.mainArea.parent
	}
}

func (g *Game) SpawnItem() {

}
