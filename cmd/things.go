package main

import (
	"math"

	"myebiten/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BULLET_RADIUS            = 4
	CHARACTER_ROTATION_SPEED = 0.05
	CHARACTER_SPEED          = 5
	CHARACTER_WIDTH          = 70
	BULLET_SPEED             = 6
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
	weapon Weapon
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
			X: c.Position.X + cos*(float64(CHARACTER_WIDTH)/2),
			Y: c.Position.Y + sin*(float64(CHARACTER_WIDTH)/2),
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
		{w.Position.X - width/2, w.Position.Y - height/2},
		{w.Position.X + width/2, w.Position.Y - height/2},
		{w.Position.X - width/2, w.Position.Y + height/2},
		{w.Position.X + width/2, w.Position.Y + height/2},
	}

	return corners
}

func (g *Game) getClosestWalls(c *Character) []*Wall {
	return nil
}

func (g *Game) getClosestWall1(b *models.Bullet) *Wall {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	i, j := getMazeCoordinates(b.Position)
	nodeCenter := getSceneCoordinates(i, j)

	if i >= len(g.Maze)-1 || i < 0 || j >= len(g.Maze[0])-1 || j < 0 {
		return nil
	}

	// Here top and bottom mean how these directions appear on the screen
	// meaning, that distToTop actually measures the distance to the wall
	// that is stored as MazeNode.bottomWall
	distToTop := b.Position.Y - (nodeCenter.Y - wh/2 + ww)
	distToRight := (nodeCenter.X + wh/2 - ww) - b.Position.X
	distToLeft := b.Position.X - (nodeCenter.X - wh/2 + ww)
	distToBottom := (nodeCenter.Y + wh/2 - ww) - b.Position.Y

	var wallToHit *Wall
	minDist := min(distToBottom, distToLeft, distToRight, distToTop)

	if minDist > b.R {
		return nil
	}

	var horizontalReflection, verticalReflection bool = false, false

	if minDist == distToBottom {
		// again, because of mismatch between how directions are logically stored
		// and how they are presented on screen bottom reflection requires MazeNode.topWall
		wallToHit = g.Maze[i][j].topWall
		if wallToHit != nil {
			// fmt.Printf("opa\n")
			// fmt.Printf("i = %d, j = %d  nodeCenter: %v pos: %v    ", i, j, nodeCenter, b.Position)
			// fmt.Printf("minDist: %f\n", minDist)
			horizontalReflection = true
		}
	}

	if minDist == distToTop {
		// see previous comments in this function
		wallToHit = g.Maze[i][j].bottomWall
		if wallToHit != nil {
			// fmt.Printf("lulw\n")
			// fmt.Printf("i = %d, j = %d  nodeCenter: %v pos: %v    ", i, j, nodeCenter, b.Position)
			// fmt.Printf("minDist: %f\n", minDist)
			horizontalReflection = true
		}
	}

	if minDist == distToLeft {
		wallToHit = g.Maze[i][j].leftWall
		if wallToHit != nil {
			verticalReflection = true
		}
	}

	if minDist == distToRight {
		wallToHit = g.Maze[i][j].rightWall
		if wallToHit != nil {
			verticalReflection = true
		}
	}

	if verticalReflection {
		cosine := math.Abs(b.Speed.X) / b.Speed.Length()
		l := minDist / cosine
		L := l * (b.R/minDist - 1)
		t := L / b.Speed.Length()

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.X = -b.Speed.X

	} else if horizontalReflection {
		cosine := math.Abs(b.Speed.Y) / b.Speed.Length()
		l := minDist / cosine
		L := l * (b.R/minDist - 1)
		t := L / b.Speed.Length()

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.Y = -b.Speed.Y
	}

	return nil
}

func (g *Game) DetectCharacterToWallCollision(c *Character) {
	closestWalls := g.getClosestWalls(c)
	for _, w := range closestWalls {
		c.detectWallCollision(*w)
	}
}

func (g *Game) DetectBulletToWallCollision(b *models.Bullet) {
	closestWall := g.getClosestWall1(b)
	if closestWall == nil {
		return
	}

	//b.Hit(&b.GameObject, closestWall.Hitbox, &closestWall.GameObject)
	// isCollision, isHorizontal := false, false

	// if int(b.position.x)%WALL_HEIGHT <= WALL_WIDTH {
	// 	wallToCollide := Wall{
	// 		X:          uint16(math.Floor(b.position.x / WALL_HEIGHT)),
	// 		Y:          uint16(math.Floor(b.position.y / WALL_HEIGHT)),
	// 		Horizontal: false,
	// 	}

	// 	_, ok := g.Walls[wallToCollide]

	// 	if ok {
	// 		isCollision, isHorizontal = true, false
	// 		//end of wall
	// 		yInWall := math.Mod(b.position.y, WALL_HEIGHT)
	// 		ySpeed := math.Abs(dy)
	// 		if yInWall <= ySpeed && dy > 0 || WALL_HEIGHT-yInWall <= ySpeed && dy < 0 {
	// 			isHorizontal = true
	// 		}
	// 	}
	// }

	// if int(b.position.y)%WALL_HEIGHT <= WALL_WIDTH {
	// 	wallToCollide := Wall{
	// 		X:          uint16(math.Floor(b.position.x / WALL_HEIGHT)),
	// 		Y:          uint16(math.Floor(b.position.y / WALL_HEIGHT)),
	// 		Horizontal: true,
	// 	}

	// 	_, ok := g.Walls[wallToCollide]

	// 	if ok {
	// 		isCollision, isHorizontal = true, true

	// 		//end of wall
	// 		xInWall := math.Mod(b.position.x, WALL_HEIGHT)
	// 		xSpeed := math.Abs(dx)
	// 		if xInWall <= xSpeed && dx > 0 || WALL_HEIGHT-xInWall <= xSpeed && dx < 0 {
	// 			isHorizontal = false
	// 		}
	// 	}
	// }

	// if isCollision {
	// 	if isHorizontal {
	// 		b.rotation = -b.rotation
	// 	}

	// 	b.rotation = math.Remainder(math.Pi-b.rotation, 2*math.Pi)
	// }
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
}
