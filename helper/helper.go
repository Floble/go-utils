package machinelearning

import (
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
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

func SumAlongColumn(m *mat.Dense) *mat.Dense {
	r, c := m.Dims()
	result := mat.NewDense(1, c, nil)

	for j := 0; j < c; j++ {
		tmp := 0.0
		for i := 0; i < r; i++ {
			tmp += m.At(i, j)
		}
		result.Set(1, j, tmp)
	}

	return result
}