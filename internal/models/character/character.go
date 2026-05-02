package character

import (
	"math"

	"myebiten/internal/models"
	"myebiten/internal/weapons"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CHARACTER_ROTATION_SPEED = 0.01
	CHARACTER_SPEED          = 0.9
	CHARACTER_WIDTH          = 60
)

type Weapon interface {
	Shoot(origin models.Vector2D, rotation float64)
}

type Character struct {
	models.GameObject
	hitbox models.RectangleHitbox
	sprite models.ImageSprite
	Input  models.Input
	weapon Weapon
}

func (c *Character) Draw(drawingArea *models.DrawingArea) {
	c.sprite.Draw(c.Position.X, c.Position.Y, c.Rotation, drawingArea)
}

func (c *Character) SetWeapon(weapon Weapon) {
	if weapon == nil {
		return
	}

	c.weapon = weapon
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
		sin, cos := math.Sincos(c.Rotation)
		origin := models.Vector2D{
			X: c.Position.X + cos*(float64(CHARACTER_WIDTH)/2+weapons.DEFAULT_GUN_BULLET_RADIUS),
			Y: c.Position.Y + sin*(float64(CHARACTER_WIDTH)/2+weapons.DEFAULT_GUN_BULLET_RADIUS),
		}

		c.weapon.Shoot(origin, c.Rotation)
	}
}

func (c *Character) getCorners() []models.Vector2D {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []models.Vector2D{
		models.RotatePoint(c.Position.X-hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		models.RotatePoint(c.Position.X+hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		models.RotatePoint(c.Position.X+hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
		models.RotatePoint(c.Position.X-hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
	}
}

func (c *Character) DetectWallCollision(wall models.Wall) bool {
	charCorners := c.getCorners()

	// Углы стены (осевой прямоугольник)
	wallCorners := wall.GetCorners()

	// Получаем оси для проверки
	axes := append(models.GetAxes(charCorners), models.GetAxes(wallCorners)...)

	// SAT: проверяем проекции на все оси
	for _, axis := range axes {
		minA, maxA := models.ProjectPolygon(axis, charCorners)
		minB, maxB := models.ProjectPolygon(axis, wallCorners)
		if !Overlap(minA, maxA, minB, maxB) {
			// Нашли разделяющую ось — значит, нет пересечения
			return false
		}
	}

	// Нет разделяющей оси — пересекаются
	return true
}

func (c *Character) DetectBulletToCharacterCollision(b *models.Bullet) (isCollision bool) {
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

func CreateCharacter(id int, charImage *ebiten.Image, weapon Weapon, controlSettings models.ControlSettings) Character {
	return Character{
		GameObject: models.GameObject{ID: id},

		hitbox: models.RectangleHitbox{H: float64(CHARACTER_WIDTH), W: float64(CHARACTER_WIDTH)},
		sprite: models.ImageSprite{Image: charImage},
		weapon: weapon,
		Input: models.Input{
			ControlSettings: controlSettings,
		},
	}
}

func Overlap(minA, maxA, minB, maxB float64) bool {
	return !(maxA < minB || maxB < minA)
}
