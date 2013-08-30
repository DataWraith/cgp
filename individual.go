package cgp

import (
	"math"
	"math/rand"
)

type Gene struct {
	Function    int
	Constant    float64
	Connections []int
}

func (g *Gene) Mutate(options *CGP) {
}

type Individual struct {
	Genes   []Gene
	Outputs []int
	Options *CGP
	Fitness float64
}

func NewIndividual(cgp *CGP) (ind Individual) {
	ind.Options = cgp
	ind.Fitness = math.Inf(1)
	ind.Genes = make([]Gene, cgp.NumGenes)
	ind.Outputs = make([]int, cgp.NumOutputs)

	for i := range ind.Genes {
		ind.Genes[i].Function = rand.Intn(len(cgp.FunctionList))
		ind.Genes[i].Constant = cgp.RandConst()
		ind.Genes[i].Connections = make([]int, cgp.MaxArity)
		for j := range ind.Genes[i].Connections {
			ind.Genes[i].Connections[j] = rand.Intn(cgp.NumInputs + i)
		}
	}

	for i := range ind.Outputs {
		ind.Outputs[i] = rand.Intn(cgp.NumInputs + cgp.NumGenes)
	}

	return
}

func (ind Individual) Mutate() (mutant Individual) {
	// Copy the parent individual
	mutant.Fitness = math.Inf(1)
	mutant.Options = ind.Options
	mutant.Genes = make([]Gene, ind.Options.NumGenes)
	mutant.Outputs = make([]int, ind.Options.NumOutputs)
	copy(mutant.Genes, ind.Genes)
	copy(mutant.Outputs, ind.Outputs)

	numMutations := ind.Options.MutationRate * float64(ind.Options.NumGenes+ind.Options.NumOutputs)
	if numMutations < 1 {
		numMutations = 1
	}

	for numMutations > 0 {
		toMutate := rand.Intn(mutant.Options.NumGenes + mutant.Options.NumOutputs)

		if toMutate < mutant.Options.NumGenes {
			mutant.Genes[toMutate].Mutate(mutant.Options)
		} else {
			mutant.Outputs[toMutate-mutant.Options.NumGenes] =
				rand.Intn(mutant.Options.NumInputs + mutant.Options.NumGenes)
		}

		numMutations--
	}

	return
}

func (ind Individual) Run(input []float64) (output []float64) {
	return []float64{0, 0, 0}
}
