package machinelearning

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
	"gonum.org/v1/gonum/floats"
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

func RandomFloat(min, max float64) float64 {
	return (max - min) * rand.Float64() + min
}

func SumAlongColumn(m *mat.Dense) *mat.Dense {
	r, c := m.Dims()
	result := mat.NewDense(1, c, nil)

	for j := 0; j < c; j++ {
		tmp := 0.0
		for i := 0; i < r - 1; i++ {
			tmp += m.At(i, j)
		}
		result.Set(0, j, tmp)
	}

	return result
}

func SumOfSquaredResiduals(output, y *mat.Dense) *mat.Dense {
	oLayerError := new(mat.Dense)
	oLayerError.Sub(y, output)
	oLayerError.Scale(-2.0, oLayerError)

	return oLayerError
}

func CrossEntropy(output, y *mat.Dense) *mat.Dense {
	oLayerError := new(mat.Dense)
	oLayerError.Sub(output, y)

	return oLayerError
}

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func DSigmoid(x float64) float64 {
	return Sigmoid(x) * (1.0 - Sigmoid(x))
}

func RelU(x float64) float64 {
	return math.Max(0.0, x)
}

func DRelU(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return 1
	}
}

func NoArg(x *mat.Dense) *mat.Dense {
	return x
}

func ArgMax(x *mat.Dense) *mat.Dense {
	r, c := x.Dims()

	for i := 0; i < r; i++ {
		tmp := make([]float64, c)
		mat.Row(tmp, i, x)
		maxIdx := floats.MaxIdx(tmp)
		
		for j := 0; j < c; j++ {
			if j == maxIdx {
				tmp[j] = 1
			} else {
				tmp[j] = 0
			}
		}

		x.SetRow(i, tmp)
	}

	return x
}

func SoftMax(x *mat.Dense) *mat.Dense {
	r, c := x.Dims()

	result := mat.DenseCopyOf(x)
	result.Apply(func(_, _ int, n float64) float64 { return math.Exp(n) }, result)

	for i := 0; i < r; i++ {
		tmp := make([]float64, c)
		mat.Row(tmp, i, result)
		sum := floats.Sum(tmp)

		for j := 0; j < c; j++ {
			result.Set(i, j, result.At(i, j) / sum)
		}
	}

	return result
}

func ReadCSV(path string, fieldsPerRecord int, resultColumns []int, inputNeurons, outputNeurons int) (*mat.Dense, *mat.Dense) {
	f, err := os.Open(path)
    if err != nil {
		fmt.Println(err)
    }
    defer f.Close()

    reader := csv.NewReader(f)
    reader.FieldsPerRecord = fieldsPerRecord

    data, err := reader.ReadAll()
    if err != nil {
		fmt.Println(err)
    }

    inputData := make([]float64, inputNeurons * len(data))
    resultData := make([]float64, outputNeurons * len(data))

    var iInput int
    var iResult int

    for i, line := range data {
		if i == 0 {
			continue
        }

        for k, v := range line {
			vParsed, err := strconv.ParseFloat(v, 64)
        	if err != nil {
				fmt.Println(err)
            }
			
			tmp := false
			for _, c := range resultColumns {
				if k == c {
					resultData[iResult] = vParsed
					iResult++
					tmp = true
					break
				}
			}
			if tmp {
				continue
			}

            inputData[iInput] = vParsed
            iInput++
        }
    }

	input := mat.NewDense(len(data), inputNeurons, inputData)
	result := mat.NewDense(len(data), outputNeurons, resultData)

	return input, result
}

func Accuracy(prediction, result *mat.Dense) float64 {
	acc := 0

	nums, _ := result.Dims()

	for i := 0; i < nums; i++ {
		if floats.MaxIdx(mat.Row(nil, i, result)) == floats.MaxIdx(mat.Row(nil, i, prediction)) {
			acc++
		}
	}

	return float64(acc) / float64(nums)
}