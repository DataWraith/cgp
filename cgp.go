package cgp

import (
	"math"
)

type CGPFunction func([]float64) float64
type EvalFunction func(Individual) float64
type RndConstFunction func() float64

type CGP struct {
	PopSize      int
	NumGenes     int
	MutationRate float64
	NumInputs    int
	NumOutputs   int
	MaxArity     int
	FunctionList []CGPFunction
	RandConst    RndConstFunction
	Evaluator    EvalFunction
	Population   []Individual
}

func New(popSize int, numGenes int, mutationRate float64, numInputs int, numOutputs int, maxArity int, functionList []CGPFunction, randomConstant RndConstFunction, evaluator EvalFunction) *CGP {

	if popSize < 2 {
		panic("Population size must be at least 2.")
	}
	if numGenes < 0 {
		panic("numGenes can't be negative")
	}
	if mutationRate < 0 || mutationRate > 1 {
		panic("Mutation rate must be between 0 and 1.")
	}
	if numInputs < 0 {
		panic("numInputs can't be negative")
	}
	if numOutputs < 1 {
		panic("At least one output is necessary.")
	}
	if maxArity < 0 {
		panic("maxArity can't be negative")
	}
	if len(functionList) == 0 {
		panic("At least one function must be provided.")
	}

	result := &CGP{
		PopSize:      popSize,
		NumGenes:     numGenes,
		MutationRate: mutationRate,
		NumInputs:    numInputs,
		NumOutputs:   numOutputs,
		MaxArity:     maxArity,
		FunctionList: functionList,
		RandConst:    randomConstant,
		Evaluator:    evaluator,
		Population:   make([]Individual, 1, popSize),
	}

	result.Population[0] = NewIndividual(result)

	return result
}

func (cgp *CGP) RunGeneration() {
	// Create offspring
	cgp.Population = cgp.Population[0:1]
	for i := 1; i < cgp.PopSize; i++ {
		cgp.Population = append(cgp.Population, cgp.Population[0].Mutate())
	}

	// Evaluate offspring
	// TODO: Parallelize this
	for i := 1; i < cgp.PopSize; i++ {
		cgp.Population[i].Fitness = cgp.Evaluator(cgp.Population[i])
	}

	// Replace parent with best offspring
	bestFitness := math.Inf(1)
	bestIndividual := 0
	for i := 1; i < cgp.PopSize; i++ {
		if cgp.Population[i].Fitness < bestFitness {
			bestFitness = cgp.Population[i].Fitness
			bestIndividual = i
		}
	}

	if bestFitness <= cgp.Population[0].Fitness {
		cgp.Population[0] = cgp.Population[bestIndividual]
	}
}
