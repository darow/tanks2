package main

import (
	"fmt"
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

func (b *Bullet) getShifts() (float64, float64) {
	sin, cos := math.Sincos(b.rotation)
	dx := cos * BULLET_SPEED
	dy := sin * BULLET_SPEED

	return dx, dy
}

func (b *Bullet) processBulletRotation(isCollision, isHorizontal bool) {
	if isCollision {
		if isHorizontal {
			b.rotation = -b.rotation
		}

		b.rotation = math.Remainder(math.Pi-b.rotation, 2*math.Pi)
	}
}

type Character struct {
	GameObject
	hitbox Hitbox
	sprite Sprite
	input  Input
	// CurrentWidth uint
}

type ControlSettings struct {
	rotateRightButton  ebiten.Key
	rotateLeftButton   ebiten.Key
	moveForwardButton  ebiten.Key
	moveBackwardButton ebiten.Key
	shootButton        ebiten.Key
}

func (c *Character) Update(bullets []*Bullet, walls map[Wall]struct{}) {
	// oldX, oldY, oldRotation := c.position.x, c.position.y, c.rotation

	if c.input.RotateRight {
		c.rotation += CHARACTER_ROTATION_SPEED
	}

	if c.input.RotateLeft {
		c.rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.input.MoveForward {
		sin, cos := math.Sincos(c.rotation)
		c.position.x += cos * CHARACTER_SPEED
		c.position.y += sin * CHARACTER_SPEED
	}

	if c.input.MoveBackward {
		sin, cos := math.Sincos(c.rotation)
		c.position.x -= cos * CHARACTER_SPEED * 5 / 6
		c.position.y -= sin * CHARACTER_SPEED * 5 / 6
	}

	// for w := range walls {
	// 	if c.detectWallCollision(w) {
	// 		fmt.Printf("opa")
	// 		c.position.x, c.position.y, c.rotation = oldX, oldY, oldRotation
	// 		break
	// 	}
	// }

	if c.input.Shoot {
		fmt.Printf("haha")
		c.input.Shoot = false
		sin, cos := math.Sincos(c.rotation)
		x := c.position.x + cos*(float64(CHARACTER_WIDTH)/2)
		y := c.position.y + sin*(float64(CHARACTER_WIDTH)/2)

		for _, bullet := range bullets {
			if !bullet.active {
				fmt.Printf("woah")
				bullet.position.x = x
				bullet.position.y = y
				bullet.rotation = c.rotation
				bullet.active = true
				break
			}
		}
	}

}

func (c *Character) getCorners() []Point {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []Point{
		rotatePoint(c.position.x-hw, c.position.y-hh, c.position.x, c.position.y, c.rotation),
		rotatePoint(c.position.x+hw, c.position.y-hh, c.position.x, c.position.y, c.rotation),
		rotatePoint(c.position.x+hw, c.position.y+hh, c.position.x, c.position.y, c.rotation),
		rotatePoint(c.position.x-hw, c.position.y+hh, c.position.x, c.position.y, c.rotation),
	}
}

func (c *Character) detectWallCollisionOld(w Wall) bool {
	// 1. Центр стены
	halfWallW := float64(WALL_WIDTH) / 2
	halfWallH := float64(WALL_HEIGHT) / 2
	if w.rotation == 0.0 {
		halfWallW, halfWallH = halfWallH, halfWallW
	}

	wallCenterX := w.position.x
	wallCenterY := w.position.y

	// 2. Смещение стены относительно персонажа
	dx := wallCenterX - c.position.x
	dy := wallCenterY - c.position.y

	// 3. Переводим в систему координат персонажа (учитывая поворот)
	xLocal := dx*math.Cos(c.rotation) + dy*math.Sin(c.rotation)
	yLocal := -dx*math.Sin(c.rotation) + dy*math.Cos(c.rotation)

	// 4. Переводим размеры
	halfCharW := float64(CHARACTER_WIDTH) / 2
	halfCharH := float64(CHARACTER_WIDTH) / 2

	// 5. Находим смещения по осям
	// Разница между центрами в локальной системе координат
	deltaX := math.Abs(xLocal)
	deltaY := math.Abs(yLocal)

	// 6. Проверяем пересечение по осям
	return deltaX <= (halfCharW+halfWallW) && deltaY <= (halfCharH+halfWallH)
}

type Wall struct {
	GameObject
	hitbox Hitbox
	sprite Sprite
}

func (w *Wall) GetCorners() []Point {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.rotation == 0.0 {
		height, width = width, height
	}

	corners := []Point{
		{w.position.x - width/2, w.position.y - height/2},
		{w.position.x + width/2, w.position.y - height/2},
		{w.position.x - width/2, w.position.y + height/2},
		{w.position.x + width/2, w.position.y + height/2},
	}

	return corners
}

func (g *Game) ProcessBullet(b *Bullet) {
	dx, dy := b.getShifts()
	b.position.x += dx
	b.position.y += dy

	isCollision, isHorizontal := false, false // g.DetectBulletToWallCollision(b, dx, dy)

	b.processBulletRotation(isCollision, isHorizontal)
}

// func (g *Game) DetectBulletToWallCollision(b Bullet, dx, dy float64) (isCollision, isHorizontal bool) {
// 	if int(b.position.x)%WALL_HEIGHT <= WALL_WIDTH {
// 		wallToCollide := Wall{
// 			X:          uint16(math.Floor(b.position.x / WALL_HEIGHT)),
// 			Y:          uint16(math.Floor(b.position.y / WALL_HEIGHT)),
// 			Horizontal: false,
// 		}

// 		_, ok := g.Walls[wallToCollide]

// 		if ok {
// 			isCollision, isHorizontal = true, false
// 			//end of wall
// 			yInWall := math.Mod(b.position.y, WALL_HEIGHT)
// 			ySpeed := math.Abs(dy)
// 			if yInWall <= ySpeed && dy > 0 || WALL_HEIGHT-yInWall <= ySpeed && dy < 0 {
// 				isHorizontal = true
// 			}

// 			return
// 		}
// 	}

// 	if int(b.position.y)%WALL_HEIGHT <= WALL_WIDTH {
// 		wallToCollide := Wall{
// 			X:          uint16(math.Floor(b.position.x / WALL_HEIGHT)),
// 			Y:          uint16(math.Floor(b.position.y / WALL_HEIGHT)),
// 			Horizontal: true,
// 		}

// 		_, ok := g.Walls[wallToCollide]

// 		if ok {
// 			isCollision, isHorizontal = true, true

// 			//end of wall
// 			xInWall := math.Mod(b.position.x, WALL_HEIGHT)
// 			xSpeed := math.Abs(dx)
// 			if xInWall <= xSpeed && dx > 0 || WALL_HEIGHT-xInWall <= xSpeed && dx < 0 {
// 				isHorizontal = false
// 			}

// 			return
// 		}
// 	}

// 	return
// }

func (g *Game) DetectBulletCharacterCollision(b *Bullet, c *Character) (isCollision bool) {
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
