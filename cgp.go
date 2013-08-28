package cgp

import (
	"math"
)

type CGPFunction func([]float64) float64
type EvalFunction func(Individual) float64
type RndConstFunction func() float64

type CGP struct {
	PopSize      uint
	MaxGenes     uint
	MutationRate float64
	NumInputs    uint
	NumOutputs   uint
	MaxArity     uint
	FunctionList []CGPFunction
	RandConst    RndConstFunction
	Evaluator    EvalFunction
	Population   []Individual
}

func New(popSize uint, maxGenes uint, mutationRate float64, numInputs uint, numOutputs uint, maxArity uint, functionList []CGPFunction, randomConstant RndConstFunction, evaluator EvalFunction) *CGP {

	if popSize < 2 {
		panic("Population size must be at least 2.")
	}

	if numOutputs == 0 {
		panic("At least one output is necessary.")
	}

	if mutationRate < 0 || mutationRate > 1 {
		panic("Mutation rate must be between 0 and 1.")
	}

	if len(functionList) == 0 {
		panic("At least one function must be provided.")
	}

	result := &CGP{
		PopSize:      popSize,
		MaxGenes:     maxGenes,
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
	for i := uint(1); i < cgp.PopSize; i++ {
		cgp.Population = append(cgp.Population, cgp.Population[0].Mutate())
	}

	// Evaluate offspring
	// TODO: Parallelize this
	for i := 1; uint(i) < cgp.PopSize; i++ {
		cgp.Population[i].Fitness = cgp.Evaluator(cgp.Population[i])
	}

	// Replace parent with best offspring
	bestFitness := math.Inf(1)
	bestIndividual := uint(0)
	for i := uint(1); i < cgp.PopSize; i++ {
		if cgp.Population[i].Fitness < bestFitness {
			bestFitness = cgp.Population[i].Fitness
			bestIndividual = i
		}
	}

	if bestFitness <= cgp.Population[0].Fitness {
		cgp.Population[0] = cgp.Population[bestIndividual]
	}
}
