package machinelearning

// Kudos to https://datadan.io/blog/neural-net-with-go

import (
	"time"
	"math/rand"
	"gonum.org/v1/gonum/mat"
	helper "go-utils/helper"
)

type NeuralNetworkConfig struct {
	inputNeurons int
	hiddenNeurons int
	outputNeurons int
	maxSteps int
	//miniBatchSize int
	learningRate float64
	//dLossFunction func(x *mat.Dense, y *mat.Dense, w *mat.Dense, b *mat.Dense) *mat.Dense
	activationFunction func(x float64) float64
	dActivationFunction func(x float64) float64
}

type NeuralNetwork struct {
	config *NeuralNetworkConfig
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOutput *mat.Dense
	bOutput *mat.Dense
}

func NewNeuralNetworkConfig(inputNeurons, hiddenNeurons, outputNeurons, maxSteps int, learningRate float64, activationFunction func(x float64) float64, dActivationFunction func(x float64) float64) *NeuralNetworkConfig {
	config := new(NeuralNetworkConfig)
	config.inputNeurons = inputNeurons
	config.hiddenNeurons = hiddenNeurons
	config.outputNeurons = outputNeurons
	config.maxSteps = maxSteps
	//config.miniBatchSize = miniBatchSize
	config.learningRate = learningRate
	//config.dLossFunction = dLossFunction
	config.activationFunction = activationFunction
	config.dActivationFunction = dActivationFunction

	return config
}

func NewNeuralNetwork(config *NeuralNetworkConfig) *NeuralNetwork {
	nn := new(NeuralNetwork)
	nn.config = config
	nn.wHidden = mat.NewDense(nn.config.inputNeurons, nn.config.hiddenNeurons, nil)
	nn.bHidden = mat.NewDense(1, nn.config.hiddenNeurons, nil)
	nn.wOutput = mat.NewDense(nn.config.hiddenNeurons, nn.config.outputNeurons, nil)
	nn.bOutput = mat.NewDense(1, nn.config.outputNeurons, nil)

	return nn
}

func (nn *NeuralNetwork) init() {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	
	for i := 0; i < nn.config.inputNeurons; i++ {
		for j := 0; j < nn.config.hiddenNeurons; j++ {
			nn.wHidden.Set(i, j, randGen.Float64())
		}
	}

	for i := 0; i < nn.config.hiddenNeurons; i++ {
		for j := 0; j < nn.config.outputNeurons; j++ {
			nn.wHidden.Set(i, j, randGen.Float64())
		}
	}

	for j := 0; j < nn.config.hiddenNeurons; j++ {
		nn.bHidden.Set(0, j, randGen.Float64())
	}

	for j := 0; j < nn.config.outputNeurons; j++ {
		nn.bHidden.Set(0, j, randGen.Float64())
	}
}

func (nn *NeuralNetwork) feedForward(x *mat.Dense) (*mat.Dense, *mat.Dense, *mat.Dense, *mat.Dense) {
	addBHidden := func(_, j int, n float64) float64 { return n + nn.bHidden.At(0, j)}
	addBOutput := func(_, j int, n float64) float64 { return n + nn.bOutput.At(0, j)}
	applyActivationFunction := func(_, _ int, n float64) float64 { return nn.config.activationFunction(n) }

	hLayerInput := new(mat.Dense)
	hLayerInput.Mul(x, nn.wHidden)
	hLayerInput.Apply(addBHidden, hLayerInput)
	hLayerOutput := new(mat.Dense)
	hLayerOutput.Apply(applyActivationFunction, hLayerInput)

	oLayerInput := new(mat.Dense)
	oLayerInput.Mul(hLayerOutput, nn.wOutput)
	oLayerInput.Apply(addBOutput, oLayerInput)
	
	output := new(mat.Dense)
	output.Apply(applyActivationFunction, oLayerInput)

	return oLayerInput, hLayerInput, hLayerOutput, output
}

func (nn *NeuralNetwork) backPropagation(x, y, output, oLayerInput, hLayerInput, hLayerOutput *mat.Dense) {
	/* dSSR/dBout = dSSR/doutput * doutput/doLayerInput * doLayerInput/dBout
	dSSR/dBout = -2 (labels - output) * sigmoid(output) * (1 - sigmoid(output)) * 1

	dSSR/dWOutput = dSSR/dOutput * dOutput/doLayerInput * doLayerInput/dWOutput
	dSSR/dWOutput = - 2 (labels - output) * sigmoid(output) * (1 - sigmoid(output)) * hLayerOutput

	dSSR/dbHidden = dSSR/dOutput * dOutput/doLayerInput * doLayerInput/dhLayerOutput * dhLayerOutput/dhLayerInput * dhLayerInput/dbHidden
	dSSR/dbHidden = -2 (labels - output) * sigmoid(output) * (1 - sigmoid(output)) * wOutput * sigmoid(hLayerOutput) * (1 - sigmoid(hLayerOutput)) * 1

	dSSR/dwHidden = dSSR/dOutput * dOutput/doLayerInput * doLayerInput/dhLayerOutput * dhLayerOutput/dhLayerInput * dhLayerInput/dwHidden
	dSSR/dwHidden = -2 (labels - output) * sigmoid(output) * (1 - sigmoid(output)) * wOutput * sigmoid(hLayerOutput) * (1 - sigmoid(hLayerOutput)) * input
	
	output = sigmoid(oLayerInput)
	oLayerInput = hLayerOutput * wOutput + bOutput
	hLayerOutput = sigmoid(hLayerInput)
	hLayerInput = input * wHidden + bHidden */
	
	applyDActivationFunction := func(_, _ int, n float64) float64 { return nn.config.dActivationFunction(n) }

	oLayerError := new(mat.Dense)
	oLayerError.Sub(y, output)
	oLayerError.Scale(-2.0, oLayerError)

	dOutput := new(mat.Dense)
	dOutput.Apply(applyDActivationFunction, output)

	dHLayer := new(mat.Dense)
	dHLayer.Apply(applyDActivationFunction, hLayerOutput)

	dBOut := new(mat.Dense)
	dBOut.MulElem(oLayerError, dOutput)

	dBHidden := new(mat.Dense)
	dBHidden.Mul(dBOut, nn.wOutput.T())
	dBHidden.MulElem(dBHidden, dHLayer)

	newBOut := helper.SumAlongColumn(dBOut)
	newBOut.Scale(nn.config.learningRate, newBOut)
	nn.bOutput.Sub(nn.bOutput, newBOut)

	newBHidden := helper.SumAlongColumn(dBHidden)
	newBHidden.Scale(nn.config.learningRate, newBHidden)
	nn.bHidden.Sub(nn.bHidden, newBHidden)

	dWOut := new(mat.Dense)
	dWOut.Mul(hLayerOutput.T(), dBOut)
	dWOut.Scale(nn.config.learningRate, dWOut)
	nn.wOutput.Sub(nn.wOutput, dWOut)

	dWHidden := new(mat.Dense)
	dWHidden.Mul(x.T(), dBHidden)
	dWHidden.Scale(nn.config.learningRate, dWHidden)
	nn.wHidden.Sub(nn.wHidden, dWHidden)
}

func (nn *NeuralNetwork) Train(x, y *mat.Dense) {
	nn.init()

	for i := 0; i < nn.config.maxSteps; i++ {
		oLayerInput, hLayerInput, hLayerOutput, output := nn.feedForward(x)
		nn.backPropagation(x, y, output, oLayerInput, hLayerInput, hLayerOutput)
	}
}

func (nn *NeuralNetwork) Predict(x *mat.Dense) *mat.Dense {
	addBHidden := func(_, j int, n float64) float64 { return n + nn.bHidden.At(0, j)}
	addBOutput := func(_, j int, n float64) float64 { return n + nn.bOutput.At(0, j)}
	applyActivationFunction := func(_, _ int, n float64) float64 { return nn.config.activationFunction(n) }

	hLayerInput := new(mat.Dense)
	hLayerInput.Mul(x, nn.wHidden)
	hLayerInput.Apply(addBHidden, hLayerInput)
	hLayerOutput := new(mat.Dense)
	hLayerOutput.Apply(applyActivationFunction, hLayerInput)

	oLayerInput := new(mat.Dense)
	oLayerInput.Mul(hLayerOutput, nn.wOutput)
	oLayerInput.Apply(addBOutput, oLayerInput)

	output := new(mat.Dense)
	output.Apply(applyActivationFunction, oLayerInput)

	return output
}