package main

import (
	"math"
)

type Vector2D struct {
	x, y float64
}

func rotatePoint(px, py, cx, cy, angle float64) Vector2D {
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

	return Vector2D{xnew, ynew}
}

func dot(a, b Vector2D) float64 {
	return a.x*b.x + a.y*b.y
}

func getAxes(points []Vector2D) []Vector2D {
	axes := []Vector2D{}
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		edge := Vector2D{p2.x - p1.x, p2.y - p1.y}
		// Нормаль
		axis := Vector2D{-edge.y, edge.x}
		// Нормализуем
		length := math.Hypot(axis.x, axis.y)
		axis.x /= length
		axis.y /= length
		axes = append(axes, axis)
	}
	return axes
}

func projectPolygon(axis Vector2D, points []Vector2D) (float64, float64) {
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

func squareDistance(v Vector2D, w Vector2D) float64 {
	return (v.x-w.x)*(v.x-w.x) + (v.y-w.y)*(v.y-w.y)
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
