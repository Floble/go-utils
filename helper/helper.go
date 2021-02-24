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
		for i := 0; i < r; i++ {
			tmp += m.At(i, j)
		}
		result.Set(0, j, tmp)
	}

	return result
}

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func DSigmoid(x float64) float64 {
	return Sigmoid(x) * (1.0 - Sigmoid(x))
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