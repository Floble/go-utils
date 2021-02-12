package machinelearning

import (
	"gonum.org/v1/gonum/mat"
	helper "go-utils/helper"
)

type NeuralNetworkConfig struct {
	nums int
	inputNeurons int
	hiddenNeurons int
	outputNeurons int
	maxSteps int
	miniBatchSize int
	learningRate float64
	dLossFunction func(x *mat.Dense, y *mat.Dense, w *mat.Dense, b *mat.Dense) *mat.Dense
	activationFunction func(x float64) float64
	dActivationFunction func(x float64) float64
}

type NeuralNetwork struct {
	config *NeuralNetworkConfig
	input *mat.Dense
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOutput *mat.Dense
	bOutput *mat.Dense
	output *mat.Dense
}

func NewNeuralNetworkConfig(nums int, inputNeurons, hiddenNeurons, outputNeurons, maxSteps, miniBatchSize int, learningRate float64, dLossFunction func(x *mat.Dense, y *mat.Dense, w *mat.Dense, b *mat.Dense) *mat.Dense, activationFunction func(x float64) float64, dActivationFunction func(x float64) float64) *NeuralNetworkConfig {
	config := new(NeuralNetworkConfig)
	config.nums = nums
	config.inputNeurons = inputNeurons
	config.hiddenNeurons = hiddenNeurons
	config.outputNeurons = outputNeurons
	config.maxSteps = maxSteps
	config.miniBatchSize = miniBatchSize
	config.learningRate = learningRate
	config.dLossFunction = dLossFunction
	config.activationFunction = activationFunction
	config.dActivationFunction = dActivationFunction

	return config
}

func NewNeuralNetwork(config *NeuralNetworkConfig, data []float64) *NeuralNetwork {
	nn := new(NeuralNetwork)
	nn.config = config
	nn.input = mat.NewDense(nn.config.nums, nn.config.inputNeurons, data)
	nn.wHidden = mat.NewDense(nn.config.inputNeurons, nn.config.hiddenNeurons, nil)
	nn.bHidden = mat.NewDense(1, nn.config.hiddenNeurons, nil)
	nn.wOutput = mat.NewDense(nn.config.hiddenNeurons, nn.config.outputNeurons, nil)
	nn.bOutput = mat.NewDense(1, nn.config.outputNeurons, nil)

	return nn
}

func (nn *NeuralNetwork) init() {
	for i := 1; i <= nn.config.inputNeurons; i++ {
		for j := 1; j <= nn.config.hiddenNeurons; j++ {
			nn.wHidden.Set(i, j, helper.RandomFloat(1))
		}
	}

	for i := 1; i <= nn.config.hiddenNeurons; i++ {
		for j := 1; j <= nn.config.outputNeurons; j++ {
			nn.wHidden.Set(i, j, helper.RandomFloat(1))
		}
	}

	for j := 1; j <= nn.config.hiddenNeurons; j++ {
		nn.bHidden.Set(1, j, 0)
	}

	for j := 1; j <= nn.config.outputNeurons; j++ {
		nn.bHidden.Set(1, j, 0)
	}
}

func (nn *NeuralNetwork) feedForward() {
	addBHidden := func(_, j int, n float64) float64 { return n + nn.bHidden.At(0, j)}
	addBOutput := func(_, j int, n float64) float64 { return n + nn.bOutput.At(0, j)}
	applyActivationFunction := func(_, _ int, n float64) float64 { return nn.config.activationFunction(n) }

	hLayerInput := new(mat.Dense)
	hLayerInput.Mul(nn.input, nn.wHidden)
	hLayerInput.Apply(addBHidden, hLayerInput)
	hLayerOutput := new(mat.Dense)
	hLayerOutput.Apply(applyActivationFunction, hLayerInput)

	oLayerInput := new(mat.Dense)
	oLayerInput.Mul(hLayerOutput, nn.wOutput)
	oLayerInput.Apply(addBOutput, oLayerInput)
	nn.output = new(mat.Dense)
	nn.output.Apply(applyActivationFunction, oLayerInput)
}

func (nn *NeuralNetwork) backPropagation(y, oLayerInput, hLayerInput, hLayerOutput *mat.Dense) {
	applyDActivationFunction := func(_, _ int, n float64) float64 { return nn.config.dActivationFunction(n) }

	oLayerError := new(mat.Dense)
	oLayerError.Sub(y, nn.output)
	oLayerError.Scale(-2.0, oLayerError)

	dOutput := new(mat.Dense)
	dOutput.Apply(applyDActivationFunction, nn.output)

	dBOut := new(mat.Dense)
	dBOut.MulElem(oLayerError, dOutput)

	dWOut := new(mat.Dense)
	dWOut.Mul(hLayerOutput.T(), dBOut)
	dWOut.Scale(nn.config.learningRate, dWOut)
	nn.wOutput.Sub(nn.wOutput, dWOut)





	dOLayerInput := new(mat.Dense)
	dOLayerInput.Apply(applyDActivationFunction, oLayerInput)

	dHLayerInput := new(mat.Dense)
	dHLayerInput.Apply(applyDActivationFunction, hLayerInput)

	oLayerError := new(mat.Dense)
	oLayerError.Sub(y, nn.output)
	oLayerError.Scale(-2.0, oLayerError)

	dBOutput := new(mat.Dense)
	dBOutput.MulElem(oLayerError, dOLayerInput)

	dWOutput := new(mat.Dense)
	dWOutput.Mul(dBOutput, hLayerOutput)


	dSSR/dBout = dSSR/doutput * doutput/doLayerInput * doLayerInput/dBout
	dSSR/dBout = -2 (labels - output) * sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * 1

	dSSR/dWOutput = dSSR/dOutput * dOutput/doLayerInput * doLayerInput/dWOutput
	dSSR/dWOutput = - 2 (labels - output) * sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * hLayerOutput

	dSSR/dbHidden = dSSR/dOutput * dOutput/doLayerInput * doLayerInput/dhLayerOutput * dhLayerOutput/dhLayerInput * dhLayerInput/dbHidden
	dSSR/dbHidden = -2 (labels - output) * sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * wOutput * sigmoid(hLayerInput) * (1 - sigmoid(hLayerInput)) * 1

	dSSR/dwHidden = dSSR/dOutput * dOutput/doLayerInput * doLayerInput/dhLayerOutput * dhLayerOutput/dhLayerInput * dhLayerInput/dwHidden
	dSSR/dwHidden = -2 (labels - output) * sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * wOutput * sigmoid(hLayerInput) * (1 - sigmoid(hLayerInput)) * input
	
	output = sigmoid(oLayerInput)
	oLayerInput = hLayerOutput * wOutput + bOutput
	hLayerOutput = sigmoid(hLayerInput)
	hLayerInput = input * wHidden + bHidden
}