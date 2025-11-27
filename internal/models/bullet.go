package models

const BULLET_RADIUS = 4

type Bullet struct {
	GameObject
	Sprite CircleSprite
	R      float64
}

func (b *Bullet) Draw(drawingArea *DrawingArea) {
	b.Sprite.Draw(b.Position.X, b.Position.Y, drawingArea)
}

func CreateBullet(r int) *Bullet {
	return &Bullet{
		Sprite: CircleSprite{R: float64(r)},
		R:      float64(r),
	}
}
