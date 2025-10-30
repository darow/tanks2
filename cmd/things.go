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
	WALL_HEIGHT              = 170 // equal to cell size in labyrinth
	WALL_WIDTH               = 10
)

var TILE_ID_SEQUENCE = 0

type Bullet struct {
	R float64
	GameObject
	Hitbox `json:"-"`
	Sprite `json:"-"`
}

type Wall struct {
	GameObject
	Hitbox `json:"-"`
	Sprite `json:"-"`
}

type Character struct {
	GameObject
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
		origin := Vector2D{
			X: c.Position.X + cos*(float64(CHARACTER_WIDTH)/2),
			Y: c.Position.Y + sin*(float64(CHARACTER_WIDTH)/2),
		}

		c.weapon.Shoot(origin, c.Rotation)
	}

}

func (c *Character) getCorners() []Vector2D {
	// Получаем точки углов персонажа (с учётом поворота)
	hw := float64(CHARACTER_WIDTH) / 2
	hh := float64(CHARACTER_WIDTH) / 2

	return []Vector2D{
		rotatePoint(c.Position.X-hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		rotatePoint(c.Position.X+hw, c.Position.Y-hh, c.Position.X, c.Position.Y, c.Rotation),
		rotatePoint(c.Position.X+hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
		rotatePoint(c.Position.X-hw, c.Position.Y+hh, c.Position.X, c.Position.Y, c.Rotation),
	}
}

func (w *Wall) GetCorners() []Vector2D {
	var height, width float64 = WALL_HEIGHT, WALL_WIDTH
	if w.Rotation == 0.0 {
		height, width = width, height
	}

	corners := []Vector2D{
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

func (g *Game) getClosestWall1(b *Bullet) *Wall {
	wh := float64(WALL_HEIGHT)
	ww := float64(WALL_WIDTH)

	a := b.Position.X / (wh - ww)
	ai := int(math.Floor(a))

	c := b.Position.Y / (wh - ww)
	ci := int(math.Floor(c))

	nodeCenter := getSceneCoordinates(ci+1, ai+1)

	if ci >= len(g.Maze)-2 || ci < 0 || ai >= len(g.Maze[0])-2 || ai < 0 {
		return nil
	}

	distToTop := (nodeCenter.Y + wh/2 - ww) - b.Position.Y
	distToRight := (nodeCenter.Y + wh/2 - ww) - b.Position.X
	distToBottom := b.Position.Y - (nodeCenter.Y - wh/2 + ww)
	distToLeft := b.Position.X - (nodeCenter.X - wh/2 + ww)

	var wallToHit *Wall
	minDist := min(distToBottom, distToLeft, distToRight, distToTop)

	if minDist > b.R {
		return nil
	}

	// yeaaaaahh, this needs to be perfromed for 3 other walls and also
	// additionally for some other
	// it's very hard to figure out the good solution ;(
	if minDist == distToBottom {
		wallToHit = g.Maze[ci+1][ai+1].bottomWall

		cosine := b.Speed.X / b.Speed.length()
		l := minDist / cosine
		t := l * (b.R/minDist - 1)

		b.Position.X -= t * b.Speed.X
		b.Position.Y -= t * b.Speed.Y

		b.Speed.Y = -b.Speed.Y
	}

	// check intersections with protrusions (fuck)

	return wallToHit
}

func lineIntersect(p1, p2, position Vector2D, speed Vector2D) (Vector2D, bool) {
	// Line segment AB represented as p1 to p2
	// Line represented as position + t * speed

	// Direction vectors for the segments
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	ldx := speed.X
	ldy := speed.Y

	denom := dx*ldy - dy*ldx

	if math.Abs(denom) < 1e-6 {
		// Lines are parallel or collinear
		return Vector2D{}, false
	}

	deltaX := position.X - p1.X
	deltaY := position.Y - p1.Y

	t := (deltaX*ldy - deltaY*ldx) / denom
	u := (deltaX*dy - deltaY*dx) / denom

	// Check if the intersection is within the segment
	if t >= 0 && t <= 1 && u >= 0 {
		// Calculate intersection point
		intersection := Vector2D{
			X: position.X + u*ldx,
			Y: position.Y + u*ldy,
		}
		return intersection, true
	}

	return Vector2D{}, false
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
