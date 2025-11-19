package models

import (
	"math"
)

func RotatePoint(px, py, cx, cy, angle float64) Vector2D {
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

	return Vector2D{X: xnew, Y: ynew}
}

func dot(a, b Vector2D) float64 {
	return a.X*b.X + a.Y*b.Y
}

func getAxes(points []Vector2D) []Vector2D {
	axes := []Vector2D{}
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		edge := Vector2D{X: p2.X - p1.X, Y: p2.Y - p1.Y}
		// Нормаль
		axis := Vector2D{X: -edge.Y, Y: edge.X}
		// Нормализуем
		length := math.Hypot(axis.X, axis.Y)
		axis.X /= length
		axis.Y /= length
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
	return (v.X-w.X)*(v.X-w.X) + (v.Y-w.Y)*(v.Y-w.Y)
}

func overlap(minA, maxA, minB, maxB float64) bool {
	return !(maxA < minB || maxB < minA)
}
