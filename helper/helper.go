package machinelearning

import (
	"math/rand"
	"time"
)

func Float64(f float64, s error) float64 {
	return f
}

func RandomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	i := max - 1 + 1
	i = rand.Intn(i) + 1

	return i
}

func RandomFloat(max int) float64 {
	return float64(RandomInt(max))
}