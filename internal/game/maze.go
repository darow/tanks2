package game

import (
	"log"
	"math"
	"math/rand"

	"myebiten/internal/models"
)

const (
	MAX_BOARD_HEIGHT = 7
	MAX_BOARD_WIDTH  = 12

	MIN_BOARD_HEIGHT = 3
	MIN_BOARD_WIDTH  = 3

	WALL_HEIGHT = 170
	WALL_WIDTH  = 10
)

type Coordinates struct {
	i, j int
}

// Struct for one square of a maze, consists of 4 bools
// each encoding whether there's a passage in corresponding direction
// which confusingly means that if e.g. up is true, then topWall is nil
// (this has already caused some confusion so probably it will warrant a rework in the future)
type MazeNode struct {
	up    bool
	down  bool
	right bool
	left  bool

	topWall, bottomWall, rightWall, leftWall *models.Wall
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

func getMazeCoordinates(pos models.Vector2D) (int, int) {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	x := (pos.X - ww/2) / (wh - ww)
	j := int(math.Floor(x))

	y := (pos.Y - ww/2) / (wh - ww)
	i := int(math.Floor(y))

	return i + 1, j + 1
}

func getSceneCoordinates(i, j int) models.Vector2D {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	return models.Vector2D{X: float64(j-1)*(wh-ww) + wh/2, Y: float64(i-1)*(wh-ww) + wh/2}
}

func (g *Game) SetupLevel() (int, int, []models.Wall) {
	h := rand.Intn(MAX_BOARD_HEIGHT-MIN_BOARD_HEIGHT) + MIN_BOARD_HEIGHT
	w := rand.Intn(MAX_BOARD_WIDTH-MIN_BOARD_WIDTH) + MIN_BOARD_WIDTH

	walls := g.CreateMaze(h, w)
	g.SetDrawingSettings(h, w)
	g.SetCharacters(h, w)

	return h, w, walls
}

func (g *Game) SetCharacters(h, w int) {
	spawnPlaces := []models.Vector2D{}
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

		char.Rotation = math.Pi / 2

		char.Speed.X = 0
		char.Speed.Y = 0

		i++
	}
}

func (g *Game) CreateMaze(h, w int) []models.Wall {
	g.Walls = make([]models.Wall, 0)

	g.Maze = createMaze(h, w)
	g.Walls = buildMaze(g.Maze, g.Walls)

	return g.Walls
}

func (g *Game) SetDrawingSettings(h, w int) {
	mainArea := g.activeScene.GetArea(MAIN_PLAYING_AREA_ID)
	if mainArea == nil {
		log.Fatal("Main playing area is not set")
	}

	areaHeight := mainArea.Height
	areaWidth := mainArea.Width

	mazeHeight := float64(h*(WALL_HEIGHT-WALL_WIDTH) + WALL_WIDTH)
	mazeWidth := float64(w*(WALL_HEIGHT-WALL_WIDTH) + WALL_WIDTH)

	scalingFactor := min(areaHeight/mazeHeight, areaWidth/mazeWidth)

	mazeHeight *= scalingFactor
	mazeWidth *= scalingFactor

	newDrawingSettings := models.DrawingSettings{
		Offset: models.Vector2D{X: (areaWidth - mazeWidth) / 2, Y: (areaHeight - mazeHeight) / 2},
		Scale:  scalingFactor,
	}

	mazeArea := mainArea.NewArea(mazeHeight, mazeWidth, newDrawingSettings)
	g.activeScene.AddDrawingArea(MAZE_AREA_ID, mazeArea)

	for _, bullet := range g.Bullets {
		g.activeScene.AddObject(bullet, MAZE_AREA_ID)
	}

	for _, char := range g.Characters {
		g.activeScene.AddObject(char, MAZE_AREA_ID)
	}

	for _, wall := range g.Walls {
		g.activeScene.AddObject(&wall, MAZE_AREA_ID)
	}
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

	// fill in missing connections on all nodes for consistency
	for i := 1; i < N+1; i++ {
		for j := 1; j < M+1; j++ {
			mazeNodes[i][j].up = mazeNodes[i][j].up || mazeNodes[i+1][j].down
			mazeNodes[i][j].right = mazeNodes[i][j].right || mazeNodes[i][j+1].left
			mazeNodes[i][j].down = mazeNodes[i][j].down || mazeNodes[i-1][j].up
			mazeNodes[i][j].left = mazeNodes[i][j].left || mazeNodes[i][j-1].right
		}
	}

	return mazeNodes
}

func buildMaze(mazeNodes [][]MazeNode, walls []models.Wall) []models.Wall {
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
				w := models.CreateWall(
					models.Vector2D{X: nodeCenter.X, Y: nodeCenter.Y - (wh-ww)/2},
					ww,
					wh,
					false,
				)
				w.SetActive(true)

				currentNode.bottomWall = &w
				downNode.topWall = &w
				walls = append(walls, w)
			}

			if verticalWall {
				w := models.CreateWall(
					models.Vector2D{X: nodeCenter.X - (wh-ww)/2, Y: nodeCenter.Y},
					ww,
					wh,
					true,
				)
				w.SetActive(true)

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
		bullet.SetActive(false)
	}

	for _, char := range g.Characters {
		char.SetActive(true)
		char.Input.Reset()
	}

	mainArea := g.activeScene.GetArea(MAIN_PLAYING_AREA_ID)
	mainArea.Children = nil
}

func (g *Game) SpawnItem() {

}
