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

type Individual struct {
	Genes   []Gene
	Options *CGP
	Fitness float64
}

func NewIndividual(cgp *CGP) (ind Individual) {
	ind.Options = cgp
	ind.Fitness = math.Inf(1)
	ind.Genes = make([]Gene, cgp.NumGenes)

	for i := range ind.Genes {
		ind.Genes[i].Function = rand.Intn(len(cgp.FunctionList))
		ind.Genes[i].Constant = cgp.RandConst()
		ind.Genes[i].Connections = make([]int, cgp.MaxArity)
		for j := range ind.Genes[i].Connections {
			ind.Genes[i].Connections[j] = rand.Intn(cgp.NumInputs + i)
		}
	}

	return
}

func (ind Individual) Mutate() Individual {
	return NewIndividual(ind.Options)
}

func (ind Individual) Run(input []float64) (output []float64) {
	return []float64{0, 0, 0}
}
