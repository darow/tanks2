package models

import "math"

type Wall struct {
	GameObject
	Hitbox RectangleHitbox `json:"-"`
	Sprite RectangleSprite `json:"-"`
}

func (w *Wall) Draw(drawingArea *DrawingArea) {
	w.Sprite.Draw(w.Position.X, w.Position.Y, w.Rotation, drawingArea)
}

func (w *Wall) GetCorners() []Vector2D {
	var height, width float64 = w.Hitbox.H, w.Hitbox.W
	if w.Rotation == 0.0 {
		height, width = width, height
	}

	corners := []Vector2D{
		{X: w.Position.X - width/2, Y: w.Position.Y - height/2},
		{X: w.Position.X + width/2, Y: w.Position.Y - height/2},
		{X: w.Position.X - width/2, Y: w.Position.Y + height/2},
		{X: w.Position.X + width/2, Y: w.Position.Y + height/2},
	}

	return corners
}

func CreateWall(center Vector2D, w, h float64, vertical bool) Wall {
	rotation := 0.0
	if vertical {
		rotation = math.Pi / 2
	}

	return Wall{
		GameObject: GameObject{
			Position: center,
			Rotation: rotation,
		},
		Hitbox: RectangleHitbox{H: h, W: w},
		Sprite: RectangleSprite{H: h, W: w},
	}
}
