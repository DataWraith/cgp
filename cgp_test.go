package cgp

import "testing"

func TestReverseInputs(t *testing.T) {
	// Simple test that evolves a function that reverses three inputs

	// First, we set up our parameters:
	const (
		popSize      = 5    // The population size. One parent plus four children.
		maxGenes     = 10   // The maximum number of functions in the genome
		mutationRate = 0.01 // The mutation rate
		numInputs    = 3    // The number of input values
		numOutputs   = 3    // The number of output values
		maxArity     = 2    // The maximum arity of the functions in the functionList
	)

	// We pass in a list of functions that can be used in the genome. Since
	// this is a toy example, we use two no-op functions that don't do
	// anything but pass one of the inputs through.
	//
	// The functions take an array of float64 values for input. The first
	// value is the constant that evolved for the function, the others come
	// from the maxArity inputs to the function.
	functionList := []CGPFunction{
		// pass through input A
		func(input []float64) float64 {
			return input[1]
		},

		// pass through input B
		func(input []float64) float64 {
			return input[2]
		},
	}

	// The evaluator punishes every mistake with +1 fitness (0 is perfect
	// fitness). We are looking for a function that reverses the three
	// inputs 1, 2, 3 into the three outputs 3, 2, 1
	evaluator := func(ind Individual) float64 {
		fitness := 0.0
		outputs := ind.Run([]float64{1, 2, 3})
		if outputs[0] != 3 {
			fitness += 1
		}
		if outputs[1] != 2 {
			fitness += 1
		}
		if outputs[2] != 1 {
			fitness += 1
		}
		return fitness
	}

	// This simple example does not use constants, so it doesn't matter what
	// this function returns. In an actual example it should return a random
	// constant for a function to use. For example, if you are evolving
	// images, it might return an integer value between 0.0 and 255.0 to use
	// as RGB value.
	randomConstant := func() float64 {
		return 0.0
	}

	// Initialize CGP
	gp := New(popSize, maxGenes, mutationRate, numInputs, numOutputs, maxArity, functionList, randomConstant, evaluator)

	// Population[0] is the parent, which is the most fit individual. We
	// loop until we've found a perfect solution (fitness 0)
	for gp.Population[0].Fitness > 0 {
		gp.RunGeneration()
	}

	t.Log("CGP successfully evolved input reversal")
}
