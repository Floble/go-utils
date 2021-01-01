package machinelearning

import (
	"math"
)

type Point struct {
	x, y float64
}

type GradientDescent struct {
	data []*Point
	learningRate float64
	maxSteps int
}

func NewPoint(x, y float64) *Point {
	p := new(Point)
	p.x = x
	p.y = y

	return p
}

func NewGradientDescent(data []*Point, learningRate float64, maxSteps int) *GradientDescent {
	gd := new(GradientDescent)
	gd.data = data
	gd.learningRate = learningRate
	gd.maxSteps = maxSteps

	return gd
}

func (gd *GradientDescent) LinearRun(sloap, intercept float64) (float64, float64) {
	stepSizeSloap := math.MaxFloat64
	stepSizeIntercept := math.MaxFloat64
	i := 1

	for (stepSizeSloap >= 0.001 && stepSizeIntercept >= 0.001) || (i <= gd.maxSteps) {
		sumSloap := 0.0
		sumIntercept := 0.0
		
		for _, p := range gd.data {
			sumSloap += -2 * p.x * (p.y - (intercept + sloap * p.x))
			sumIntercept += -2 * (p.y - (intercept + sloap * p.x))
		}

		stepSizeSloap = sumSloap * gd.learningRate
		stepSizeIntercept = sumIntercept * gd.learningRate

		sloap = sloap - stepSizeSloap
		intercept = intercept - stepSizeIntercept

		i++
	}

	return sloap, intercept
}