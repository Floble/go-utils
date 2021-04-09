package machinelearning

import (
	"math"
	vector "github.com/atedja/go-vector"
	helper "go-utils/helper"
)

type StochasticGradientDescent struct {
	data []*Point
	derivedLossFunction func(w vector.Vector, x *Point) vector.Vector
	learningRate float64
	maxSteps int
	miniBatchSize int
}

func NewStochasticGradientDescent(data []*Point, derivedLossFunction func(w vector.Vector, x *Point) vector.Vector, learningRate float64, maxSteps int, miniBatchSize int) *StochasticGradientDescent {
	sgd := new(StochasticGradientDescent)
	sgd.data = data
	sgd.derivedLossFunction = derivedLossFunction
	sgd.learningRate = learningRate
	sgd.maxSteps = maxSteps
	sgd.miniBatchSize = miniBatchSize

	return sgd
}

func (sgd *StochasticGradientDescent) gradient(w vector.Vector, d int) vector.Vector {
	gradient := vector.New(d)
	gradient.Zero()

	for i := 0; i < sgd.miniBatchSize; i++ {
		j := helper.RandomInt(len(sgd.data) - 1)
		gradient = vector.Add(gradient, sgd.derivedLossFunction(w, sgd.data[j]))
	}

	gradient.Scale(1.0 / float64(sgd.miniBatchSize))
	return gradient
}

func (sgd *StochasticGradientDescent) Run(d int) vector.Vector {
	w := vector.New(d)
	w.Zero()

	for i := 1; i <= sgd.maxSteps; i++ {
		gradient := sgd.gradient(w, d)
		sgd.learningRate = 0.1 / math.Sqrt(float64(i))
		gradient.Scale(sgd.learningRate)
		w = vector.Subtract(w, gradient)
	}

	return w
}