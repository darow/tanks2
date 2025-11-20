package game

import (
	"fmt"
	"math"

	"myebiten/internal/models"
)

const (
	ROOT_AREA_ID         = "root_area"
	MAIN_PLAYING_AREA_ID = "main_playing_area"
	MAZE_AREA_ID         = "maze_area"
	UI_AREA1_ID          = "ui_area_1"
	SCORE_AREA_ID        = "score_area"
	SCORE_AREA_1_ID      = "score_area_1"
	SCORE_AREA_2_ID      = "score_area_2"
	SCORE_AREA_3_ID      = "score_area_3"
	SCORE_AREA_4_ID      = "score_area_4"
)

const (
	MAIN_SCENE_ID = 1
)

var TILE_ID_SEQUENCE = 0

func (g *Game) updateScores(id int) {
	g.CharactersScores[id]++
	g.scoreUITexts[id].SetText(fmt.Sprintf("Player %d: %d", id+1, g.CharactersScores[id]))
}

func (g *Game) getClosestWalls(c *models.Character) []*models.Wall {
	// yes, this is shit, I see it too, dw it will all change
	i, j := getMazeCoordinates(c.Position)
	k := 0
	if g.Maze[i][j].topWall != nil {
		wallsToCheck[k] = g.Maze[i][j].topWall
		k++
	}
	if g.Maze[i][j].bottomWall != nil {
		wallsToCheck[k] = g.Maze[i][j].bottomWall
		k++
	}
	if g.Maze[i][j].leftWall != nil {
		wallsToCheck[k] = g.Maze[i][j].leftWall
		k++
	}
	if g.Maze[i][j].rightWall != nil {
		wallsToCheck[k] = g.Maze[i][j].rightWall
		k++
	}

	if g.Maze[i-1][j].leftWall != nil {
		wallsToCheck[k] = g.Maze[i-1][j].leftWall
		k++
	}
	if g.Maze[i-1][j].rightWall != nil {
		wallsToCheck[k] = g.Maze[i-1][j].rightWall
		k++
	}

	if g.Maze[i+1][j].leftWall != nil {
		wallsToCheck[k] = g.Maze[i+1][j].leftWall
		k++
	}
	if g.Maze[i+1][j].rightWall != nil {
		wallsToCheck[k] = g.Maze[i+1][j].rightWall
		k++
	}

	if g.Maze[i][j-1].topWall != nil {
		wallsToCheck[k] = g.Maze[i][j-1].topWall
		k++
	}
	if g.Maze[i][j-1].bottomWall != nil {
		wallsToCheck[k] = g.Maze[i][j-1].bottomWall
		k++
	}

	if g.Maze[i][j+1].topWall != nil {
		wallsToCheck[k] = g.Maze[i][j+1].topWall
		k++
	}
	if g.Maze[i][j+1].bottomWall != nil {
		wallsToCheck[k] = g.Maze[i][j+1].bottomWall
		k++
	}

	return wallsToCheck[:k]
}

func (g *Game) DetectCharacterToWallCollision(c *models.Character) {
	closestWalls := g.getClosestWalls(c)
	for _, w := range closestWalls {
		isCollide := c.DetectWallCollision(*w)
		if isCollide {
			c.MoveBack()
		}
	}
}

func (g *Game) DetectBulletToWallCollision(b *models.Bullet) {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	i, j := getMazeCoordinates(b.Position)
	nodeCenter := getSceneCoordinates(i, j)

	if i >= len(g.Maze)-1 || i < 0 || j >= len(g.Maze[0])-1 || j < 0 {
		return
	}

	// Here top and bottom mean how these directions appear on the screen
	// meaning, that distToTop actually measures the distance to the wall
	// that is stored as MazeNode.bottomWall
	distToTop := b.Position.Y - (nodeCenter.Y - wh/2 + ww)
	distToRight := (nodeCenter.X + wh/2 - ww) - b.Position.X
	distToLeft := b.Position.X - (nodeCenter.X - wh/2 + ww)
	distToBottom := (nodeCenter.Y + wh/2 - ww) - b.Position.Y

	minDist := min(distToBottom, distToLeft, distToRight, distToTop)

	if minDist > b.R {
		return
	}

	var horizontalReflection, verticalReflection bool = false, false

	if minDist == distToBottom {
		// again, because of mismatch between how directions are logically stored
		// and how they are presented on screen bottom reflection requires MazeNode.topWall

		// check if the top wall is present
		horizontalReflection = !g.Maze[i][j].up

		// if close to the left check 3 corner walls, same if close to the right
		// if at least one corner wall is present, perform reflection
		verticalReflection = (distToLeft < b.R) && !(g.Maze[i][j].up && g.Maze[i+1][j].left && g.Maze[i][j-1].up) ||
			(distToRight < b.R) && !(g.Maze[i][j].up && g.Maze[i+1][j].right && g.Maze[i][j+1].up)

		// prioritize reflection of the main wall
		verticalReflection = verticalReflection && !horizontalReflection
	}

	if minDist == distToTop {
		// check comments in minDist == distToBottom block
		horizontalReflection = !g.Maze[i][j].down

		verticalReflection = (distToLeft < b.R) && !(g.Maze[i][j].down && g.Maze[i-1][j].left && g.Maze[i][j-1].down) ||
			(distToRight < b.R) && !(g.Maze[i][j].down && g.Maze[i-1][j].right && g.Maze[i][j+1].down)

		verticalReflection = verticalReflection && !horizontalReflection
	}

	if minDist == distToLeft {
		// check comments in minDist == distToBottom block
		verticalReflection = !g.Maze[i][j].left

		horizontalReflection = (distToTop < b.R) && !(g.Maze[i][j].down && g.Maze[i-1][j].left && g.Maze[i][j-1].down) ||
			(distToBottom < b.R) && !(g.Maze[i][j].up && g.Maze[i+1][j].left && g.Maze[i][j-1].up)

		horizontalReflection = horizontalReflection && !verticalReflection
	}

	if minDist == distToRight {
		// check comments in minDist == distToBottom block
		verticalReflection = !g.Maze[i][j].right

		horizontalReflection = (distToTop < b.R) && !(g.Maze[i][j].down && g.Maze[i-1][j].right && g.Maze[i][j+1].down) ||
			(distToBottom < b.R) && !(g.Maze[i][j].up && g.Maze[i+1][j].right && g.Maze[i][j+1].up)

		horizontalReflection = horizontalReflection && !verticalReflection
	}

	if verticalReflection {
		cosine := math.Abs(b.Speed.X) / b.Speed.Length()
		l := minDist / cosine
		L := b.R / cosine
		t := (L - l) / b.Speed.Length()

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.X = -b.Speed.X

	} else if horizontalReflection {
		cosine := math.Abs(b.Speed.Y) / b.Speed.Length()
		l := minDist / cosine
		L := b.R / cosine
		t := (L - l) / b.Speed.Length()

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.Y = -b.Speed.Y
	}
}

func (g *Game) SetActiveScene(sceneID int) {
	g.activeScene = g.scenes[sceneID]
}
