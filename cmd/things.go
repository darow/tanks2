package main

import (
	"math"

	"myebiten/internal/models"
	"myebiten/internal/weapons"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BULLET_RADIUS            = 4
	CHARACTER_ROTATION_SPEED = 0.05
	CHARACTER_SPEED          = 5
	CHARACTER_WIDTH          = 70
	WALL_HEIGHT              = 170
	WALL_WIDTH               = 10
)

var TILE_ID_SEQUENCE = 0

type Wall struct {
	models.GameObject
	Hitbox `json:"-"`
	Sprite `json:"-"`
}

type Character struct {
	models.GameObject
	Hitbox `json:"-"`
	Sprite `json:"-"`
	input  Input
	weapon weapons.Weapon
}

type ControlSettings struct {
	rotateRightButton  ebiten.Key
	rotateLeftButton   ebiten.Key
	moveForwardButton  ebiten.Key
	moveBackwardButton ebiten.Key
	shootButton        ebiten.Key
}

func (c *Character) ProcessInput() {
	c.Speed.X = 0.0
	c.Speed.Y = 0.0

	if c.input.RotateRight {
		c.Rotation += CHARACTER_ROTATION_SPEED
	}

	if c.input.RotateLeft {
		c.Rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.input.MoveForward {
		sin, cos := math.Sincos(c.Rotation)
		c.Speed.X = cos * CHARACTER_SPEED
		c.Speed.Y = sin * CHARACTER_SPEED
	}

	if c.input.MoveBackward {
		sin, cos := math.Sincos(c.Rotation)
		c.Speed.X = -cos * CHARACTER_SPEED * 5 / 6
		c.Speed.Y = -sin * CHARACTER_SPEED * 5 / 6
	}

	if c.input.Shoot {
		c.input.Shoot = false
		sin, cos := math.Sincos(c.Rotation)
		origin := models.Vector2D{
			X: c.Position.X + cos*(float64(CHARACTER_WIDTH)/2+BULLET_RADIUS),
			Y: c.Position.Y + sin*(float64(CHARACTER_WIDTH)/2+BULLET_RADIUS),
		}

		c.weapon.Shoot(origin, c.Rotation)
	}

}

func (c *Character) getCorners() []models.Vector2D {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []models.Vector2D{
		RotatePoint(c.Position.X-hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		RotatePoint(c.Position.X+hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		RotatePoint(c.Position.X+hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
		RotatePoint(c.Position.X-hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
	}
}

func (w *Wall) GetCorners() []models.Vector2D {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.Rotation == 0.0 {
		height, width = width, height
	}

	corners := []models.Vector2D{
		{X: w.Position.X - width/2, Y: w.Position.Y - height/2},
		{X: w.Position.X + width/2, Y: w.Position.Y - height/2},
		{X: w.Position.X - width/2, Y: w.Position.Y + height/2},
		{X: w.Position.X + width/2, Y: w.Position.Y + height/2},
	}

	return corners
}

func (g *Game) getClosestWalls(c *Character) []*Wall {
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

func (g *Game) DetectCharacterToWallCollision(c *Character) {
	closestWalls := g.getClosestWalls(c)
	for _, w := range closestWalls {
		isCollide := c.detectWallCollision(*w)
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

func (g *Game) DetectBulletToCharacterCollision(b *models.Bullet, c *Character) (isCollision bool) {
	// Сдвигаем снаряд в локальную систему координат прямоугольника
	dx := b.Position.X - c.Position.X
	dy := b.Position.Y - c.Position.Y

	sin, cos := math.Sincos(c.Rotation)

	xLocal := dx*cos + dy*sin
	yLocal := -dx*sin + dy*cos

	// Находим ближайшую точку на прямоугольнике
	halfW := float64(CHARACTER_WIDTH) / 2
	halfH := float64(CHARACTER_WIDTH) / 2

	closestX := math.Max(-halfW, math.Min(xLocal, halfW))
	closestY := math.Max(-halfH, math.Min(yLocal, halfH))

	// Вычисляем расстояние от снаряда до ближайшей точки
	dxLocal := xLocal - closestX
	dyLocal := yLocal - closestY
	distanceSq := dxLocal*dxLocal + dyLocal*dyLocal

	// Сравниваем в квадратах, чтоб не вычислять корень из distanceSq
	return distanceSq <= BULLET_RADIUS*BULLET_RADIUS
}

func (c *Character) Copy(c2 *Character) {
	if c2 == nil {
		c.Position.X = 99999
		c.Position.Y = 99999
		return
	}

	c.Position.X = c2.Position.X
	c.Position.Y = c2.Position.Y
	c.Rotation = c2.Rotation

	c.Active = c2.Active
}
