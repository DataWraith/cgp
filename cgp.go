package cgp

import (
	"math"
	"sync"
)

type CGPFunction func([]float64) float64
type EvalFunction func(Individual) float64
type RndConstFunction func() float64

type CGPOptions struct {
	PopSize      int
	NumGenes     int
	MutationRate float64
	NumInputs    int
	NumOutputs   int
	MaxArity     int
	FunctionList []CGPFunction
	RandConst    RndConstFunction
	Evaluator    EvalFunction
}

type cgp struct {
	Options    CGPOptions
	Population []Individual
}

func New(options CGPOptions) *cgp {

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

	result := &cgp{
		Options:    options,
		Population: make([]Individual, 1, options.PopSize),
	}

	result.Population[0] = NewIndividual(&options)

	return result
}

func (cgp *cgp) RunGeneration() {
	// Create offspring
	cgp.Population = cgp.Population[0:1]
	for i := 1; i < cgp.Options.PopSize; i++ {
		cgp.Population = append(cgp.Population, cgp.Population[0].Mutate())
	}

	// Evaluate offspring (in parallel)
	var wg sync.WaitGroup
	for i := 1; i < cgp.Options.PopSize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cgp.Population[i].Fitness = cgp.Options.Evaluator(cgp.Population[i])
		}(i)
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
