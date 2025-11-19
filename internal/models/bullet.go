package models

const BULLET_RADIUS = 4

type Bullet struct {
	GameObject
	sprite CircleSprite
	R      float64
}

func (b *Bullet) Draw(drawingArea *DrawingArea) {
	b.sprite.Draw(b.Position.X, b.Position.Y, drawingArea)
}

func CreateBullet(r int) *Bullet {
	return &Bullet{
		sprite: CircleSprite{R: float64(r)},
		R:      float64(r),
	}
}
