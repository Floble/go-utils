package machinelearning

// Kudos to https://datadan.io/blog/neural-net-with-go
// Data for training and testing is obtained from https://github.com/dwhitena/gophernet/tree/master/data

import (
	helper "go-utils/helper"
	"math/rand"
	"time"
	"gonum.org/v1/gonum/mat"
)

type NeuralNetworkConfig struct {
	inputNeurons int
	hiddenNeurons int
	outputNeurons int
	maxSteps int
	//miniBatchSize int
	learningRate float64
	lossFunction func(output *mat.Dense, y *mat.Dense) *mat.Dense
	activationFunction func(x float64) float64
	dActivationFunction func(x float64) float64
	outputFunctionTraining func(output *mat.Dense) *mat.Dense
	outputFunctionPrediction func(output *mat.Dense) *mat.Dense
}

type NeuralNetwork struct {
	config *NeuralNetworkConfig
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOutput *mat.Dense
	bOutput *mat.Dense
}

func NewNeuralNetworkConfig(inputNeurons, hiddenNeurons, outputNeurons, maxSteps int, learningRate float64, lossFunction func(output *mat.Dense, y *mat.Dense) *mat.Dense, activationFunction func(x float64) float64, dActivationFunction func(x float64) float64, outputFunctionTraining func(output *mat.Dense) *mat.Dense, outputFunctionPrediction func(output *mat.Dense) *mat.Dense) *NeuralNetworkConfig {
	config := new(NeuralNetworkConfig)
	config.inputNeurons = inputNeurons
	config.hiddenNeurons = hiddenNeurons
	config.outputNeurons = outputNeurons
	config.maxSteps = maxSteps
	//config.miniBatchSize = miniBatchSize
	config.learningRate = learningRate
	config.lossFunction = lossFunction
	config.activationFunction = activationFunction
	config.dActivationFunction = dActivationFunction
	config.outputFunctionTraining = outputFunctionTraining
	config.outputFunctionPrediction = outputFunctionPrediction

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
			nn.wOutput.Set(i, j, randGen.Float64())
		}
	}

	for j := 0; j < nn.config.hiddenNeurons; j++ {
		nn.bHidden.Set(0, j, randGen.Float64())
	}

	for j := 0; j < nn.config.outputNeurons; j++ {
		nn.bOutput.Set(0, j, randGen.Float64())
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
	
	oRaw := new(mat.Dense)
	oRaw.Apply(applyActivationFunction, oLayerInput)
	output := nn.config.outputFunctionTraining(oRaw)

	return oLayerInput, hLayerInput, hLayerOutput, output
}

func (nn *NeuralNetwork) backPropagation(x, y, output, oLayerInput, hLayerInput, hLayerOutput *mat.Dense) {
	/* <<<<< Loss Function = SSR >>>>>

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
	hLayerInput = input * wHidden + bHidden */

	/* <<<<< Loss Function = CE >>>>>
	
	dCE/dBout = dCE/dOSoftMax * dOSoftMax/dORaw * dORaw/dOLayerInput * dOLayerInput/dBOut
	dCE/dBout = (-1/oSoftMax) * (oSoftMax * (1 - oSoftMax)) || ((-1/oSoftMax') * (-oSoftMax * oSoftMax')) * (sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * 1
	dCE/dBout = (oSoftMax - 1) || (oSoftMax) * (sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * 1

	dCE/dWOutput = dCE/dOSoftMax * dOSoftMax/dORaw * dORaw/dOLayerInput * dOLayerInput/dWOutput
	dCE/dWOutput = (-1/oSoftMax) * (oSoftMax * (1 - oSoftMax)) || ((-1/oSoftMax') * (-oSoftMax * oSoftMax')) * (sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * hLayerOutput
	dCE/dWOutput = (oSoftMax - 1) || (oSoftMax) * (sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * hLayerOutput

	dCE/dBHidden = dCE/dOSoftMax * dOSoftMax/dORaw * dORaw/dOLayerInput * dOLayerInput/dHLayerOutput * dHLayerOutput/dHLayerInput * dHLayerInput/dBHidden
	dCE/dBHidden = (oSoftMax - 1) || (oSoftMax) * (sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * wOutput * sigmoid(hLayerInput) * (1 - sigmoid(hLayerInput)) * 1

	dCE/dWHidden = dCE/dOSoftMax * dOSoftMax/dORaw * dORaw/dOLayerInput * dOLayerInput/dHLayerOutput * dHLayerOutput/dHLayerInput * dHLayerInput/dWHidden
	dCE/dWHidden = (oSoftMax - 1) || (oSoftMax) * (sigmoid(oLayerInput) * (1 - sigmoid(oLayerInput)) * wOutput * sigmoid(hLayerInput) * (1 - sigmoid(hLayerInput)) * input

	CE = -log(oSoftMax)
	oSoftMax = softmax(oRaw)
	oRaw = sigmoid(oLayerInput) 
	oLayerInput = hLayerOutput * wOutput + bOutput 
	hLayerOutput = sigmoid(hLayerInput) 
	hLayerInput = input * wHidden + bHidden */
	
	applyDActivationFunction := func(_, _ int, n float64) float64 { return nn.config.dActivationFunction(n) }

	oLayerError := nn.config.lossFunction(output, y)

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
	_, _, _, output := nn.feedForward(x)
	output = nn.config.outputFunctionPrediction(output)

	return output
}