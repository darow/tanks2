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
	WALL_HEIGHT              = 150 // equal to cell size in labyrinth
	WALL_WIDTH               = 10
)

var TILE_ID_SEQUENCE = 0

type WallsDTO struct {
	Walls []Wall `json:"walls"`
}

type Bullet struct {
	GameObject
	hitbox Hitbox
	sprite Sprite
}

type Character struct {
	GameObject
	hitbox Hitbox
	sprite Sprite
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
	c.speed.x = 0.0
	c.speed.y = 0.0

	if c.input.RotateRight {
		c.rotation += CHARACTER_ROTATION_SPEED
	}

	if c.input.RotateLeft {
		c.rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.input.MoveForward {
		sin, cos := math.Sincos(c.rotation)
		c.speed.x = cos * CHARACTER_SPEED
		c.speed.y = sin * CHARACTER_SPEED
	}

	if c.input.MoveBackward {
		sin, cos := math.Sincos(c.rotation)
		c.speed.x = -cos * CHARACTER_SPEED * 5 / 6
		c.speed.y = -sin * CHARACTER_SPEED * 5 / 6
	}

	if c.input.Shoot {
		c.input.Shoot = false
		sin, cos := math.Sincos(c.rotation)
		origin := Vector2D{
			x: c.position.x + cos*(float64(CHARACTER_WIDTH)/2),
			y: c.position.y + sin*(float64(CHARACTER_WIDTH)/2),
		}

		c.weapon.Shoot(origin, c.rotation)
	}

}

func (c *Character) getCorners() []Vector2D {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []Vector2D{
		rotatePoint(c.position.x-hw, c.position.y-hh, c.position.x, c.position.y, c.rotation),
		rotatePoint(c.position.x+hw, c.position.y-hh, c.position.x, c.position.y, c.rotation),
		rotatePoint(c.position.x+hw, c.position.y+hh, c.position.x, c.position.y, c.rotation),
		rotatePoint(c.position.x-hw, c.position.y+hh, c.position.x, c.position.y, c.rotation),
	}
}

type Wall struct {
	GameObject
	hitbox Hitbox
	sprite Sprite
}

func (w *Wall) GetCorners() []Vector2D {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.rotation == 0.0 {
		height, width = width, height
	}

	corners := []Vector2D{
		{w.position.x - width/2, w.position.y - height/2},
		{w.position.x + width/2, w.position.y - height/2},
		{w.position.x - width/2, w.position.y + height/2},
		{w.position.x + width/2, w.position.y + height/2},
	}

	return corners
}

func (g *Game) getClosestWalls(c *Character) []Wall {
	return nil
}

func (g *Game) DetectCharacterToWallCollision(c *Character) {
	closestWalls := g.getClosestWalls(c)
	for _, w := range closestWalls {
		c.detectWallCollision(w)
	}
}

func (g *Game) DetectBulletToWallCollision(b *Bullet) {
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
	dx := b.position.x - c.position.x
	dy := b.position.y - c.position.y

	sin, cos := math.Sincos(c.rotation)

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
		c.position.x = 99999
		c.position.y = 99999
		return
	}

	c.position.x = c2.position.x
	c.position.y = c2.position.y
	c.rotation = c2.rotation
}
