package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	BULLET_RADIUS            = 4
	CHARACTER_ROTATION_SPEED = 0.038
	CHARACTER_SPEED          = 3.5
	CHARACTER_WIDTH          = 40
	BULLET_SPEED             = 4.3
	WALL_HEIGHT              = 100 // equal to cell size in labyrinth
	WALL_WIDTH               = 10
)

var TILE_ID_SEQUENCE = 0

type Things struct {
	bullets map[int]Bullet
	walls   map[Wall]struct{}
}

type Bullet struct {
	id       int
	x, y     float64
	rotation float64
}

func (t *Things) getNextID() int {
	TILE_ID_SEQUENCE++

	if len(t.bullets) >= math.MaxUint16 {
		log.Fatal("can not get next id: all possible values are used. Increase id type to uint32")
	}

	if _, ok := t.bullets[TILE_ID_SEQUENCE]; ok {
		return t.getNextID()
	}

	return TILE_ID_SEQUENCE
}

func (b *Bullet) getShifts() (float64, float64) {
	sin, cos := math.Sincos(b.rotation)
	dx := cos * BULLET_SPEED
	dy := sin * BULLET_SPEED

	return dx, dy
}

func (b *Bullet) processBulletRotation(isCollision, isHorizontal bool) Bullet {
	if isCollision {
		if isHorizontal {
			b.rotation = -b.rotation

			return *b
		}

		b.rotation = math.Remainder(math.Pi-b.rotation, 2*math.Pi)
	}

	return *b
}

type Character struct {
	x, y     float64
	rotation float64

	charImg      *ebiten.Image
	currentWidth uint

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
	oldX, oldY, oldRotation := c.x, c.y, c.rotation

	if c.input.rotateRight {
		c.rotation += CHARACTER_ROTATION_SPEED
	}

	if c.input.rotateLeft {
		c.rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.input.moveForward {
		sin, cos := math.Sincos(c.rotation)
		c.x += cos * CHARACTER_SPEED
		c.y += sin * CHARACTER_SPEED
	}

	if c.input.moveBackward {
		sin, cos := math.Sincos(c.rotation)
		c.x -= cos * CHARACTER_SPEED * 5 / 6
		c.y -= sin * CHARACTER_SPEED * 5 / 6
	}

	for w := range walls {
		if c.detectWallCollision(w) {
			c.x, c.y, c.rotation = oldX, oldY, oldRotation
			break
		}
	}

	if inpututil.IsKeyJustPressed(c.input.shootButton) {
		sin, cos := math.Sincos(c.rotation)
		x := c.x + cos*(float64(c.currentWidth)/2)
		y := c.y + sin*(float64(c.currentWidth)/2)

		b := Bullet{
			x:        x,
			y:        y,
			rotation: c.rotation,
		}

		return &b
	}

	return nil
}

func (c *Character) getCorners() []Point {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(c.currentWidth) / 2
	hh := float64(c.currentWidth) / 2

	return []Point{
		rotatePoint(c.x-hw, c.y-hh, c.x, c.y, c.rotation),
		rotatePoint(c.x+hw, c.y-hh, c.x, c.y, c.rotation),
		rotatePoint(c.x+hw, c.y+hh, c.x, c.y, c.rotation),
		rotatePoint(c.x-hw, c.y+hh, c.x, c.y, c.rotation),
	}
}

func (c *Character) detectWallCollisionOld(w Wall) bool {
	// 1. Центр стены
	halfWallW := float64(WALL_WIDTH) / 2
	halfWallH := float64(WALL_HEIGHT) / 2
	if w.horizontal {
		halfWallW, halfWallH = halfWallH, halfWallW
	}

	wallCenterX := float64(w.x)*WALL_HEIGHT + halfWallW
	wallCenterY := float64(w.y)*WALL_HEIGHT + halfWallH

	// 2. Смещение стены относительно персонажа
	dx := wallCenterX - c.x
	dy := wallCenterY - c.y

	// 3. Переводим в систему координат персонажа (учитывая поворот)
	xLocal := dx*math.Cos(c.rotation) + dy*math.Sin(c.rotation)
	yLocal := -dx*math.Sin(c.rotation) + dy*math.Cos(c.rotation)

	// 4. Переводим размеры
	halfCharW := float64(c.currentWidth) / 2
	halfCharH := float64(c.currentWidth) / 2

	// 5. Находим смещения по осям
	// Разница между центрами в локальной системе координат
	deltaX := math.Abs(xLocal)
	deltaY := math.Abs(yLocal)

	// 6. Проверяем пересечение по осям
	return deltaX <= (halfCharW+halfWallW) && deltaY <= (halfCharH+halfWallH)
}

type Wall struct {
	x, y       uint16
	horizontal bool
}

// GetCenter deprecated
func (w *Wall) GetCenter() (x, y float32) {
	x = WALL_WIDTH/2 + float32(w.x)*WALL_HEIGHT
	y = WALL_HEIGHT/2 + float32(w.y)*WALL_HEIGHT
	if w.horizontal {
		x, y = y, x
	}

	return x, y
}

func (w *Wall) GetCorners() []Point {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.horizontal {
		height, width = width, height
	}

	corners := []Point{
		{float64(w.x) * WALL_HEIGHT, float64(w.y) * WALL_HEIGHT},
		{float64(w.x)*WALL_HEIGHT + width, float64(w.y) * WALL_HEIGHT},
		{float64(w.x) * WALL_HEIGHT, float64(w.y)*WALL_HEIGHT + height},
		{float64(w.x)*WALL_HEIGHT + width, float64(w.y)*WALL_HEIGHT + height},
	}

	return corners
}

func (t *Things) ProcessBullet(b Bullet) Bullet {
	dx, dy := b.getShifts()
	b.x += dx
	b.y += dy

	isCollision, isHorizontal := t.DetectBulletToWallCollision(b, dx, dy)

	b = b.processBulletRotation(isCollision, isHorizontal)

	return b
}

func (t *Things) DetectBulletToWallCollision(b Bullet, dx, dy float64) (isCollision, isHorizontal bool) {
	if int(b.x)%WALL_HEIGHT <= WALL_WIDTH {
		wallToCollide := Wall{
			x:          uint16(math.Floor(b.x / WALL_HEIGHT)),
			y:          uint16(math.Floor(b.y / WALL_HEIGHT)),
			horizontal: false,
		}

		if _, ok := t.walls[wallToCollide]; ok {
			isCollision, isHorizontal = true, false
			//end of wall
			yInWall := math.Mod(b.y, WALL_HEIGHT)
			ySpeed := math.Abs(dy)
			if yInWall <= ySpeed && dy > 0 || WALL_HEIGHT-yInWall <= ySpeed && dy < 0 {
				isHorizontal = true
			}

			return
		}
	}

	if int(b.y)%WALL_HEIGHT <= WALL_WIDTH {
		wallToCollide := Wall{
			x:          uint16(math.Floor(b.x / WALL_HEIGHT)),
			y:          uint16(math.Floor(b.y / WALL_HEIGHT)),
			horizontal: true,
		}

		if _, ok := t.walls[wallToCollide]; ok {
			isCollision, isHorizontal = true, true

			//end of wall
			xInWall := math.Mod(b.x, WALL_HEIGHT)
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
	dx := b.x - c.x
	dy := b.y - c.y

	sin, cos := math.Sincos(c.rotation)

	xLocal := dx*cos + dy*sin
	yLocal := -dx*sin + dy*cos

	// Находим ближайшую точку на прямоугольнике
	halfW := float64(c.currentWidth) / 2
	halfH := float64(c.currentWidth) / 2

	closestX := math.Max(-halfW, math.Min(xLocal, halfW))
	closestY := math.Max(-halfH, math.Min(yLocal, halfH))

	// Вычисляем расстояние от снаряда до ближайшей точки
	dxLocal := xLocal - closestX
	dyLocal := yLocal - closestY
	distanceSq := dxLocal*dxLocal + dyLocal*dyLocal

	// Сравниваем в квадратах, чтоб не вычислять корень из distanceSq
	return distanceSq <= BULLET_RADIUS*BULLET_RADIUS
}
