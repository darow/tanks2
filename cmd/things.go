package main

import (
	"math"

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

type WallsDTO struct {
	Walls []Wall `json:"walls"`
}

type Bullet struct {
	R float64
	GameObject
	Hitbox
	Sprite
}

type Wall struct {
	GameObject
	Hitbox
	Sprite
}

type Character struct {
	GameObject
	Hitbox
	Sprite
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
	c.Speed.x = 0.0
	c.Speed.y = 0.0

	if c.input.RotateRight {
		c.Rotation += CHARACTER_ROTATION_SPEED
	}

	if c.input.RotateLeft {
		c.Rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.input.MoveForward {
		sin, cos := math.Sincos(c.Rotation)
		c.Speed.x = cos * CHARACTER_SPEED
		c.Speed.y = sin * CHARACTER_SPEED
	}

	if c.input.MoveBackward {
		sin, cos := math.Sincos(c.Rotation)
		c.Speed.x = -cos * CHARACTER_SPEED * 5 / 6
		c.Speed.y = -sin * CHARACTER_SPEED * 5 / 6
	}

	if c.input.Shoot {
		c.input.Shoot = false
		sin, cos := math.Sincos(c.Rotation)
		origin := Vector2D{
			x: c.Position.x + cos*(float64(CHARACTER_WIDTH)/2),
			y: c.Position.y + sin*(float64(CHARACTER_WIDTH)/2),
		}

		c.weapon.Shoot(origin, c.Rotation)
	}

}

func (c *Character) getCorners() []Vector2D {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []Vector2D{
		rotatePoint(c.Position.x-hw, c.Position.y-hh, c.Position.x, c.Position.y, c.Rotation),
		rotatePoint(c.Position.x+hw, c.Position.y-hh, c.Position.x, c.Position.y, c.Rotation),
		rotatePoint(c.Position.x+hw, c.Position.y+hh, c.Position.x, c.Position.y, c.Rotation),
		rotatePoint(c.Position.x-hw, c.Position.y+hh, c.Position.x, c.Position.y, c.Rotation),
	}
}

func (w *Wall) GetCorners() []Vector2D {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.Rotation == 0.0 {
		height, width = width, height
	}

	corners := []Vector2D{
		{w.Position.x - width/2, w.Position.y - height/2},
		{w.Position.x + width/2, w.Position.y - height/2},
		{w.Position.x - width/2, w.Position.y + height/2},
		{w.Position.x + width/2, w.Position.y + height/2},
	}

	return corners
}

func (g *Game) getClosestWalls(c *Character) []*Wall {
	return nil
}

func (g *Game) getClosestWall1(b *Bullet) *Wall {
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
	distToTop := b.Position.y - (nodeCenter.y - wh/2 + ww)
	distToRight := (nodeCenter.x + wh/2 - ww) - b.Position.x
	distToLeft := b.Position.x - (nodeCenter.x - wh/2 + ww)
	distToBottom := (nodeCenter.y + wh/2 - ww) - b.Position.y

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
		cosine := math.Abs(b.Speed.x) / b.Speed.length()
		l := minDist / cosine
		L := l * (b.R/minDist - 1)
		t := L / b.Speed.length()

		b.Position.x -= t * b.Speed.x
		b.Position.y -= t * b.Speed.y

		b.Speed.x = -b.Speed.x

	} else if horizontalReflection {
		cosine := math.Abs(b.Speed.y) / b.Speed.length()
		l := minDist / cosine
		L := l * (b.R/minDist - 1)
		t := L / b.Speed.length()

		b.Position.x -= t * b.Speed.x
		b.Position.y -= t * b.Speed.y

		b.Speed.y = -b.Speed.y
	}

	return nil
}

func (g *Game) DetectCharacterToWallCollision(c *Character) {
	closestWalls := g.getClosestWalls(c)
	for _, w := range closestWalls {
		c.detectWallCollision(*w)
	}
}

func (g *Game) DetectBulletToWallCollision(b *Bullet) {
	closestWall := g.getClosestWall1(b)
	if closestWall == nil {
		return
	}

	b.Hit(&b.GameObject, closestWall.Hitbox, &closestWall.GameObject)
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

func (g *Game) DetectBulletToCharacterCollision(b *Bullet, c *Character) (isCollision bool) {
	// Сдвигаем снаряд в локальную систему координат прямоугольника
	dx := b.Position.x - c.Position.x
	dy := b.Position.y - c.Position.y

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
		c.Position.x = 99999
		c.Position.y = 99999
		return
	}

	c.Position.x = c2.Position.x
	c.Position.y = c2.Position.y
	c.Rotation = c2.Rotation
}
