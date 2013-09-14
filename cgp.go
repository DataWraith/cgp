// Package cgp implements Cartesian Genetic Programming in Go.
package cgp

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// A CGPFunction is a function that is usable in a Genetic Program. It takes
// zero or more parameters and outputs a single result. For example
// a CGPFunction could implement binary AND or floating point multiplication.
type CGPFunction func([]float64) float64

// The EvalFunction takes one Individual and returns its fitness value.
type EvalFunction func(Individual) float64

// RndConstFunction takes a PRNG as input and outputs a random number that is
// used as a constant in the evolved program. This allows you to set the range
// and type (integers vs. floating point) of constants used during evolution.
// For example, if you are evolving programs that create RGB images you might
// constrain the RndConstFunction to return integer values between 0 and 255.
type RndConstFunction func(rand *rand.Rand) float64

// CGPOptions is a struct describing the options of a CGP run.
type CGPOptions struct {
	PopSize      int              // Population Size
	NumGenes     int              // Number of Genes
	MutationRate float64          // Mutation Rate
	NumInputs    int              // The number of Inputs
	NumOutputs   int              // The number of Outputs
	MaxArity     int              // The maximum Arity of the CGPFunctions in FunctionList
	FunctionList []CGPFunction    // The functions used in evolution
	RandConst    RndConstFunction // The function supplying constants
	Evaluator    EvalFunction     // The evaluator that assigns a fitness to an individual
	Rand         *rand.Rand       // An instance of rand.Rand that is used throughout cgp to make runs repeatable
}

type CGP struct {
	Options        CGPOptions
	Population     []Individual
	NumEvaluations int // The number of evaluations so far
}

// New takes CGPOptions and returns a new CGP object. It panics when a necessary
// precondition is violated, e.g. when the number of genes is negative.
func New(options CGPOptions) *CGP {

	if options.PopSize < 2 {
		panic("Population size must be at least 2.")
	}
	if options.NumGenes < 0 {
		panic("NumGenes can't be negative.")
	}
	if options.MutationRate < 0 || options.MutationRate > 1 {
		panic("Mutation rate must be between 0 and 1.")
	}
	if options.NumInputs < 0 {
		panic("NumInputs can't be negative.")
	}
	if options.NumOutputs < 1 {
		panic("At least one output is necessary.")
	}
	if options.MaxArity < 0 {
		panic("MaxArity can't be negative.")
	}
	if len(options.FunctionList) == 0 {
		panic("At least one function must be provided.")
	}
	if options.RandConst == nil {
		panic("You must supply a RandConst function.")
	}
	if options.Evaluator == nil {
		panic("You must supply an Evaluator function.")
	}

	if options.Rand == nil {
		options.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	result := &CGP{
		Options:        options,
		Population:     make([]Individual, 1, options.PopSize),
		NumEvaluations: 0,
	}

	result.Population[0] = NewIndividual(&options)

	return result
}

// RunGeneration creates offspring from the current parent via mutation,
// evaluates the offspring using the CGP object's Evaluator and selects the new
// parent for the following generation.
func (cgp *CGP) RunGeneration() {
	// Create offspring
	cgp.Population = cgp.Population[0:1]
	for i := 1; i < cgp.Options.PopSize; i++ {
		cgp.Population = append(cgp.Population, cgp.Population[0].Mutate())
	}

	// Evaluate offspring (in parallel)
	var wg sync.WaitGroup
	for i := 1; i < cgp.Options.PopSize; i++ {
		// If the individual computes the same function as the parent, skip
		// evaluation and just use the parent's fitness
		if cgp.Population[i].CacheID() == cgp.Population[0].CacheID() {
			cgp.Population[i].Fitness = cgp.Population[0].Fitness
		} else {
			// Individual is different from parent, compute fitness
			wg.Add(1)
			cgp.NumEvaluations += 1
			go func(i int) {
				defer wg.Done()
				cgp.Population[i].Fitness = cgp.Options.Evaluator(cgp.Population[i])
			}(i)
		}
	}
	wg.Wait()

	// Replace parent with best offspring
	bestFitness := math.Inf(1)
	bestIndividual := 0
	for i := 1; i < cgp.Options.PopSize; i++ {
		if cgp.Population[i].Fitness < bestFitness {
			bestFitness = cgp.Population[i].Fitness
			bestIndividual = i
		}
	}

	if bestFitness <= cgp.Population[0].Fitness {
		cgp.Population[0] = cgp.Population[bestIndividual]
	}
}
