package machinelearning

import (
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