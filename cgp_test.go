package cgp_test

import (
	"fmt"
	"github.com/DataWraith/cgp"
	"math"
	"math/rand"
	"runtime"
	"testing"
)

func TestReverseInputs(t *testing.T) {
	// Simple test that evolves a function that reverses three inputs

	// First, we set up our parameters:
	options := cgp.CGPOptions{
		PopSize:      5,    // The population size. One parent plus four children.
		NumGenes:     10,   // The maximum number of functions in the genome
		MutationRate: 0.01, // The mutation rate
		NumInputs:    3,    // The number of input values
		NumOutputs:   3,    // The number of output values
		MaxArity:     2,    // The maximum arity of the functions in the FunctionList
	}

	// We pass in a list of functions that can be used in the genome. Since
	// this is a toy example, we use two no-op functions that don't do
	// anything but pass one of the inputs through.
	//
	// The functions take an array of float64 values for input. The first
	// value is the constant that evolved for the function, the others come
	// from the maxArity inputs to the function.
	options.FunctionList = []cgp.CGPFunction{
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
	options.Evaluator = func(ind cgp.Individual) float64 {
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
	options.RandConst = func(rand *rand.Rand) float64 {
		return 0.0
	}

	// Initialize CGP
	gp := cgp.New(options)

	// Population[0] is the parent, which is the most fit individual. We
	// loop until we've found a perfect solution (fitness 0)
	for gp.Population[0].Fitness > 0 {
		gp.RunGeneration()
	}

	t.Log("CGP successfully evolved input reversal")
}

func TestSymbolicRegression(t *testing.T) {
	// Try to approximate the function x**3 - 2x + 10

	// First, we set up our parameters:
	options := cgp.CGPOptions{
		PopSize:      100,  // The population size. One parent plus four children.
		NumGenes:     30,   // The maximum number of functions in the genome
		MutationRate: 0.03, // The mutation rate
		NumInputs:    1,    // The number of input values
		NumOutputs:   1,    // The number of output values
		MaxArity:     2,    // The maximum arity of the functions in the FunctionList
	}

	options.FunctionList = []cgp.CGPFunction{
		func(input []float64) float64 { return input[0] },            // Constant
		func(input []float64) float64 { return input[1] },            // Pass through A
		func(input []float64) float64 { return input[2] },            // Pass through B
		func(input []float64) float64 { return input[1] + input[2] }, // Addition
		func(input []float64) float64 { return input[1] - input[2] }, // Subtraction
		func(input []float64) float64 { return input[1] * input[2] }, // Multiplication

		// Division
		func(input []float64) float64 {
			if input[2] == 0 {
				return 0
			}
			return input[1] / input[2]
		},
	}

	// Generate random integer constants between 0 and 100
	options.RandConst = func(rand *rand.Rand) float64 {
		return float64(rand.Intn(101))
	}

	// Prepare the testcases.
	var testCases = []struct {
		in  float64
		out float64
	}{
		{0, 10},
		{0.5, 9.125},
		{1, 9},
		{10, 990},
		{-5, -105},
		{17, 4889},
		{3.14, 34.679144},
	}

	options.Evaluator = func(ind cgp.Individual) float64 {
		fitness := 0.0
		for _, tc := range testCases {
			input := []float64{tc.in}
			output := ind.Run(input)
			fitness += math.Pow(tc.out-output[0], 2)
		}
		return fitness
	}

	// Initialize CGP
	gp := cgp.New(options)

	// Make sure we're using all CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Population[0] is the parent, which is the most fit individual. We
	// loop until we've found a perfect solution (fitness 0)
	for gp.Population[0].Fitness > 0 {
		gp.RunGeneration()
		fmt.Println(gp.Population[0].Fitness)
	}

	t.Log("CGP successfully evolved x**3 - 2x + 10")
}
