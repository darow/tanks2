package models

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CHARACTER_ROTATION_SPEED = 0.05
	CHARACTER_SPEED          = 5
	CHARACTER_WIDTH          = 70
)

type Weapon interface {
	Shoot(origin Vector2D, rotation float64)
	Discharge()
}

type Character struct {
	GameObject
	hitbox RectangleHitbox
	sprite ImageSprite `json:"-"`
	Input  Input
	weapon Weapon
}

func (c *Character) Draw(drawingArea *DrawingArea) {
	c.sprite.Draw(c.Position.X, c.Position.Y, c.Rotation, drawingArea)
}

func (c *Character) ProcessInput() {
	c.Speed.X = 0.0
	c.Speed.Y = 0.0

	if c.Input.RotateRight {
		c.Rotation += CHARACTER_ROTATION_SPEED
	}

	if c.Input.RotateLeft {
		c.Rotation -= CHARACTER_ROTATION_SPEED
	}

	if c.Input.MoveForward {
		sin, cos := math.Sincos(c.Rotation)
		c.Speed.X = cos * CHARACTER_SPEED
		c.Speed.Y = sin * CHARACTER_SPEED
	}

	if c.Input.MoveBackward {
		sin, cos := math.Sincos(c.Rotation)
		c.Speed.X = -cos * CHARACTER_SPEED * 5 / 6
		c.Speed.Y = -sin * CHARACTER_SPEED * 5 / 6
	}

	if c.Input.Shoot {
		c.Input.Shoot = false
		sin, cos := math.Sincos(c.Rotation)
		origin := Vector2D{
			X: c.Position.X + cos*(float64(CHARACTER_WIDTH)/2+BULLET_RADIUS),
			Y: c.Position.Y + sin*(float64(CHARACTER_WIDTH)/2+BULLET_RADIUS),
		}

		c.weapon.Shoot(origin, c.Rotation)
	}

}

func (c *Character) getCorners() []Vector2D {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []Vector2D{
		RotatePoint(c.Position.X-hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		RotatePoint(c.Position.X+hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		RotatePoint(c.Position.X+hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
		RotatePoint(c.Position.X-hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
	}
}

func (c *Character) DetectWallCollision(wall Wall) bool {
	charCorners := c.getCorners()

	// Углы стены (осевой прямоугольник)
	wallCorners := wall.GetCorners()

	// Получаем оси для проверки
	axes := append(getAxes(charCorners), getAxes(wallCorners)...)

	// SAT: проверяем проекции на все оси
	for _, axis := range axes {
		minA, maxA := projectPolygon(axis, charCorners)
		minB, maxB := projectPolygon(axis, wallCorners)
		if !overlap(minA, maxA, minB, maxB) {
			// Нашли разделяющую ось — значит, нет пересечения
			return false
		}
	}

	// Нет разделяющей оси — пересекаются
	return true
}

func (c *Character) DetectBulletToCharacterCollision(b *Bullet) (isCollision bool) {
	// Сдвигаем снаряд в локальную систему координат прямоугольника
	dx := b.Position.X - c.Position.X
	dy := b.Position.Y - c.Position.Y

	sin, cos := math.Sincos(c.Rotation)

	xLocal := dx*cos + dy*sin
	yLocal := -dx*sin + dy*cos

	// Находим ближайшую точку на прямоугольнике
	halfW := float64(c.hitbox.W) / 2
	halfH := float64(c.hitbox.H) / 2

	closestX := math.Max(-halfW, math.Min(xLocal, halfW))
	closestY := math.Max(-halfH, math.Min(yLocal, halfH))

	// Вычисляем расстояние от снаряда до ближайшей точки
	dxLocal := xLocal - closestX
	dyLocal := yLocal - closestY
	distanceSq := dxLocal*dxLocal + dyLocal*dyLocal

	// Сравниваем в квадратах, чтоб не вычислять корень из distanceSq
	return distanceSq <= b.R*b.R
}

func CreateCharacter(id int, charImage *ebiten.Image, weapon Weapon, controlSettings ControlSettings) Character {
	return Character{
		GameObject: GameObject{ID: id},

		hitbox: RectangleHitbox{H: float64(CHARACTER_WIDTH), W: float64(CHARACTER_WIDTH)},
		sprite: ImageSprite{Image: charImage},
		weapon: weapon,
		Input: Input{
			ControlSettings: controlSettings,
		},
	}
}
