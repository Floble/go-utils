package machinelearning

import (
	helper "go-utils/helper"
	"math/rand"
	"time"
	vector "github.com/atedja/go-vector"
)

type GradientDescent struct {
	data []*Point
	derivedLossFunction func(w vector.Vector, x *Point) vector.Vector
	learningRate float64
	maxSteps int
}

func NewGradientDescent(data []*Point, derivedLossFunction func(w vector.Vector, x *Point) vector.Vector, learningRate float64, maxSteps int) *GradientDescent {
	gd := new(GradientDescent)
	gd.data = data
	gd.derivedLossFunction = derivedLossFunction
	gd.learningRate = learningRate
	gd.maxSteps = maxSteps

	return gd
}

func (gd *GradientDescent) gradient(w vector.Vector, d int) vector.Vector {
	gradient := vector.New(d)
	gradient.Zero()

	for _, p := range gd.data {
		gradient = vector.Add(gradient, gd.derivedLossFunction(w, p))
	}

	gradient.Scale(1.0 / float64(len(gd.data)))
	return gradient
}

func (gd *GradientDescent) Run(d int) vector.Vector {
	w := vector.New(d)
	w.Zero()
	i := 0

	for i <= gd.maxSteps {
		gradient := gd.gradient(w, d)
		gradient.Scale(gd.learningRate)
		w = vector.Subtract(w, gradient)
		i++
	}

	return w
}

func GenerateTestVector(d, max int) vector.Vector {
	v := vector.New(d)
	data := make([]float64, 0)

	for i := 0; i < d; i++ {
		rand := generateRandomNumber(max)
		data = append(data, float64(rand))
	}

	v.Set(data)
	return v
}

func GenerateTestData(d, max, amount int, w vector.Vector) []*Point {
	points := make([]*Point, 0)

	for i := 0; i < amount; i++ {
		x := GenerateTestVector(d, max)
		tmp := x.Clone()
		y := helper.Float64(vector.Dot(w, tmp))
		p := NewPoint(x, y)
		points = append(points, p)
	}

	return points
}

func generateRandomNumber(max int) int {
	rand.Seed(time.Now().UnixNano())
	i := max - 1 + 1
	i = rand.Intn(i) + 1

	return i
}