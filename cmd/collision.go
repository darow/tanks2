package main

import (
	"math"

	"myebiten/internal/models"
)

func RotatePoint(px, py, cx, cy, angle float64) models.Vector2D {
	s := math.Sin(angle)
	c := math.Cos(angle)

	// Переносим в начало координат
	px -= cx
	py -= cy

	// Вращаем
	xnew := px*c - py*s
	ynew := px*s + py*c

	// Возвращаем обратно
	xnew += cx
	ynew += cy

	return models.Vector2D{xnew, ynew}
}

func dot(a, b models.Vector2D) float64 {
	return a.X*b.X + a.Y*b.Y
}

func getAxes(points []models.Vector2D) []models.Vector2D {
	axes := []models.Vector2D{}
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		edge := models.Vector2D{p2.X - p1.X, p2.Y - p1.Y}
		// Нормаль
		axis := models.Vector2D{-edge.Y, edge.X}
		// Нормализуем
		length := math.Hypot(axis.X, axis.Y)
		axis.X /= length
		axis.Y /= length
		axes = append(axes, axis)
	}
	return axes
}

func projectPolygon(axis models.Vector2D, points []models.Vector2D) (float64, float64) {
	min := dot(points[0], axis)
	max := min
	for _, p := range points[1:] {
		proj := dot(p, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}
	return min, max
}

func squareDistance(v models.Vector2D, w models.Vector2D) float64 {
	return (v.X-w.X)*(v.X-w.X) + (v.Y-w.Y)*(v.Y-w.Y)
}

func overlap(minA, maxA, minB, maxB float64) bool {
	return !(maxA < minB || maxB < minA)
}

func (c *Character) detectWallCollision(wall Wall) bool {
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
