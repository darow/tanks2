package main

import (
	"math"
)

type Vec2 struct {
	x, y float64
}

func rotatePoint(px, py, cx, cy, angle float64) Vec2 {
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

	return Vec2{xnew, ynew}
}

func dot(a, b Vec2) float64 {
	return a.x*b.x + a.y*b.y
}

func getAxes(points []Vec2) []Vec2 {
	axes := []Vec2{}
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		edge := Vec2{p2.x - p1.x, p2.y - p1.y}
		// Нормаль
		axis := Vec2{-edge.y, edge.x}
		// Нормализуем
		length := math.Hypot(axis.x, axis.y)
		axis.x /= length
		axis.y /= length
		axes = append(axes, axis)
	}
	return axes
}

func projectPolygon(axis Vec2, points []Vec2) (float64, float64) {
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
