package machinelearning

import (
	"math"
	vector "github.com/atedja/go-vector"
)

type StochasticGradientDescent struct {
	data []*Point
	derivedLossFunction func(w vector.Vector, x *Point) vector.Vector
	learningRate float64
	maxSteps int
}

func NewStochasticGradientDescent(data []*Point, derivedLossFunction func(w vector.Vector, x *Point) vector.Vector, learningRate float64, maxSteps int) *StochasticGradientDescent {
	sgd := new(StochasticGradientDescent)
	sgd.data = data
	sgd.derivedLossFunction = derivedLossFunction
	sgd.learningRate = learningRate
	sgd.maxSteps = maxSteps

	return sgd
}

func (sgd *StochasticGradientDescent) Run(d int) vector.Vector {
	w := vector.New(d)
	w.Zero()

	for i := 1; i <= sgd.maxSteps; i++ {
		j := generateRandomNumber(len(sgd.data) - 1)
		gradient := sgd.derivedLossFunction(w, sgd.data[j])
		sgd.learningRate = 0.1 / math.Sqrt(float64(i))
		gradient.Scale(sgd.learningRate)
		w = vector.Subtract(w, gradient)
	}

	return w
}