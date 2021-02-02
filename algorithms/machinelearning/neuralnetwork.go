package machinelearning

import (
	"gonum.org/v1/gonum/mat"
)

type NeuralNetworkConfig struct {
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
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOutput *mat.Dense
	bOutput *mat.Dense
}

func NewNeuralNetworkConfig(inputNeurons, hiddenNeurons, outputNeurons, maxSteps, miniBatchSize int, learningRate float64, dLossFunction func(x *mat.Dense, y *mat.Dense, w *mat.Dense, b *mat.Dense) *mat.Dense, activationFunction func(x float64) float64, dActivationFunction func(x float64) float64) *NeuralNetworkConfig {
	config := new(NeuralNetworkConfig)
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

func NewNeuralNetwork(config *NeuralNetworkConfig) *NeuralNetwork {
	nn := new(NeuralNetwork)
	nn.config = config
	nn.wHidden = mat.NewDense(nn.config.inputNeurons, nn.config.hiddenNeurons, nil)
	nn.bHidden = mat.NewDense(1, nn.config.hiddenNeurons, nil)
	nn.wOutput = mat.NewDense(nn.config.hiddenNeurons, nn.config.outputNeurons, nil)
	nn.bOutput = mat.NewDense(1, nn.config.outputNeurons, nil)

	return nn
}