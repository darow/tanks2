package main

import (
	"log"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BULLET_RADIUS            = 4
	CHARACTER_ROTATION_SPEED = 0.05
	CHARACTER_SPEED          = 5
	CHARACTER_WIDTH          = 70
	BULLET_SPEED             = 6
	WALL_HEIGHT              = 200 // equal to cell size in labyrinth
	WALL_WIDTH               = 10
)

var TILE_ID_SEQUENCE = 0

type Things struct {
	Bullets map[int]Bullet
	walls   map[Wall]struct{}
	wallsMu sync.RWMutex
}

type WallsDTO struct {
	Walls []Wall `json:"walls"`
}

type Bullet struct {
	ID       int
	X, Y     float64
	Rotation float64
}

func (t *Things) getNextID() int {
	TILE_ID_SEQUENCE++

	if len(t.Bullets) >= math.MaxUint16 {
		log.Fatal("can not get next id: all possible values are used. Increase id type to uint32")
	}

	if _, ok := t.Bullets[TILE_ID_SEQUENCE]; ok {
		return t.getNextID()
	}

	return TILE_ID_SEQUENCE
}

func (b *Bullet) getShifts() (float64, float64) {
	sin, cos := math.Sincos(b.Rotation)
	dx := cos * BULLET_SPEED
	dy := sin * BULLET_SPEED

	return dx, dy
}

func (b *Bullet) processBulletRotation(isCollision, isHorizontal bool) Bullet {
	if isCollision {
		if isHorizontal {
			b.Rotation = -b.Rotation

			return *b
		}

		b.Rotation = math.Remainder(math.Pi-b.Rotation, 2*math.Pi)
	}

	return *b
}

type Character struct {
	id int

	X, Y     float64
	Rotation float64

	charImg      *ebiten.Image
	CurrentWidth uint

	input Input
}

type ControlSettings struct {
	rotateRightButton  ebiten.Key
	rotateLeftButton   ebiten.Key
	moveForwardButton  ebiten.Key
	moveBackwardButton ebiten.Key
	shootButton        ebiten.Key
}

func (c *Character) Update(walls map[Wall]struct{}) *Bullet {
	oldX, oldY, oldRotation := c.X, c.Y, c.Rotation

	if c.input.RotateRight {
		c.Rotation += CHARACTER_ROTATION_SPEED
	}

	if c.input.RotateLeft {
		c.Rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.input.MoveForward {
		sin, cos := math.Sincos(c.Rotation)
		c.X += cos * CHARACTER_SPEED
		c.Y += sin * CHARACTER_SPEED
	}

	if c.input.MoveBackward {
		sin, cos := math.Sincos(c.Rotation)
		c.X -= cos * CHARACTER_SPEED * 5 / 6
		c.Y -= sin * CHARACTER_SPEED * 5 / 6
	}

	for w := range walls {
		if c.detectWallCollision(w) {
			c.X, c.Y, c.Rotation = oldX, oldY, oldRotation
			break
		}
	}

	if c.input.Shoot {
		c.input.Shoot = false
		sin, cos := math.Sincos(c.Rotation)
		x := c.X + cos*(float64(c.CurrentWidth)/2)
		y := c.Y + sin*(float64(c.CurrentWidth)/2)

		b := Bullet{
			X:        x,
			Y:        y,
			Rotation: c.Rotation,
		}

		return &b
	}

	return nil
}

func (c *Character) getCorners() []Point {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(c.CurrentWidth) / 2
	hh := float64(c.CurrentWidth) / 2

	return []Point{
		rotatePoint(c.X-hw, c.Y-hh, c.X, c.Y, c.Rotation),
		rotatePoint(c.X+hw, c.Y-hh, c.X, c.Y, c.Rotation),
		rotatePoint(c.X+hw, c.Y+hh, c.X, c.Y, c.Rotation),
		rotatePoint(c.X-hw, c.Y+hh, c.X, c.Y, c.Rotation),
	}
}

func (c *Character) detectWallCollisionOld(w Wall) bool {
	// 1. Центр стены
	halfWallW := float64(WALL_WIDTH) / 2
	halfWallH := float64(WALL_HEIGHT) / 2
	if w.Horizontal {
		halfWallW, halfWallH = halfWallH, halfWallW
	}

	wallCenterX := float64(w.X)*WALL_HEIGHT + halfWallW
	wallCenterY := float64(w.Y)*WALL_HEIGHT + halfWallH

	// 2. Смещение стены относительно персонажа
	dx := wallCenterX - c.X
	dy := wallCenterY - c.Y

	// 3. Переводим в систему координат персонажа (учитывая поворот)
	xLocal := dx*math.Cos(c.Rotation) + dy*math.Sin(c.Rotation)
	yLocal := -dx*math.Sin(c.Rotation) + dy*math.Cos(c.Rotation)

	// 4. Переводим размеры
	halfCharW := float64(c.CurrentWidth) / 2
	halfCharH := float64(c.CurrentWidth) / 2

	// 5. Находим смещения по осям
	// Разница между центрами в локальной системе координат
	deltaX := math.Abs(xLocal)
	deltaY := math.Abs(yLocal)

	// 6. Проверяем пересечение по осям
	return deltaX <= (halfCharW+halfWallW) && deltaY <= (halfCharH+halfWallH)
}

type Wall struct {
	X, Y       uint16
	Horizontal bool
}

// GetCenter deprecated
func (w *Wall) GetCenter() (x, y float32) {
	x = WALL_WIDTH/2 + float32(w.X)*WALL_HEIGHT
	y = WALL_HEIGHT/2 + float32(w.Y)*WALL_HEIGHT
	if w.Horizontal {
		x, y = y, x
	}

	return x, y
}

func (w *Wall) GetCorners() []Point {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.Horizontal {
		height, width = width, height
	}

	corners := []Point{
		{float64(w.X) * WALL_HEIGHT, float64(w.Y) * WALL_HEIGHT},
		{float64(w.X)*WALL_HEIGHT + width, float64(w.Y) * WALL_HEIGHT},
		{float64(w.X) * WALL_HEIGHT, float64(w.Y)*WALL_HEIGHT + height},
		{float64(w.X)*WALL_HEIGHT + width, float64(w.Y)*WALL_HEIGHT + height},
	}

	return corners
}

func (t *Things) ProcessBullet(b Bullet) Bullet {
	dx, dy := b.getShifts()
	b.X += dx
	b.Y += dy

	isCollision, isHorizontal := t.DetectBulletToWallCollision(b, dx, dy)

	b = b.processBulletRotation(isCollision, isHorizontal)

	return b
}

func (t *Things) DetectBulletToWallCollision(b Bullet, dx, dy float64) (isCollision, isHorizontal bool) {
	if int(b.X)%WALL_HEIGHT <= WALL_WIDTH {
		wallToCollide := Wall{
			X:          uint16(math.Floor(b.X / WALL_HEIGHT)),
			Y:          uint16(math.Floor(b.Y / WALL_HEIGHT)),
			Horizontal: false,
		}

		t.wallsMu.RLock()
		_, ok := t.walls[wallToCollide]
		t.wallsMu.RUnlock()
		if ok {
			isCollision, isHorizontal = true, false
			//end of wall
			yInWall := math.Mod(b.Y, WALL_HEIGHT)
			ySpeed := math.Abs(dy)
			if yInWall <= ySpeed && dy > 0 || WALL_HEIGHT-yInWall <= ySpeed && dy < 0 {
				isHorizontal = true
			}

			return
		}
	}

	if int(b.Y)%WALL_HEIGHT <= WALL_WIDTH {
		wallToCollide := Wall{
			X:          uint16(math.Floor(b.X / WALL_HEIGHT)),
			Y:          uint16(math.Floor(b.Y / WALL_HEIGHT)),
			Horizontal: true,
		}

		t.wallsMu.RLock()
		_, ok := t.walls[wallToCollide]
		t.wallsMu.RUnlock()
		if ok {
			isCollision, isHorizontal = true, true

			//end of wall
			xInWall := math.Mod(b.X, WALL_HEIGHT)
			xSpeed := math.Abs(dx)
			if xInWall <= xSpeed && dx > 0 || WALL_HEIGHT-xInWall <= xSpeed && dx < 0 {
				isHorizontal = false
			}

			return
		}
	}

	return
}

func (t *Things) DetectBulletCharacterCollision(b Bullet, c *Character) (isCollision bool) {
	// Сдвигаем снаряд в локальную систему координат прямоугольника
	dx := b.X - c.X
	dy := b.Y - c.Y

	sin, cos := math.Sincos(c.Rotation)

	xLocal := dx*cos + dy*sin
	yLocal := -dx*sin + dy*cos

	// Находим ближайшую точку на прямоугольнике
	halfW := float64(c.CurrentWidth) / 2
	halfH := float64(c.CurrentWidth) / 2

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
		c.X = 99999
		c.Y = 99999
		return
	}

	c.X = c2.X
	c.Y = c2.Y
	c.Rotation = c2.Rotation
}
