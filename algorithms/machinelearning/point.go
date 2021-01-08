package machinelearning

import (
	vector "github.com/atedja/go-vector"
)

type Point struct {
	X vector.Vector
	Y float64
}

func NewPoint(x vector.Vector, y float64) *Point {
	p := new(Point)
	p.X = x
	p.Y = y

	return p
}