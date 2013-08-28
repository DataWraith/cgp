package cgp

import (
	"math"
)

type Individual struct {
	Options *CGP
	Fitness float64
}

func NewIndividual(cgp *CGP) Individual {
	return Individual{
		Options: cgp,
		Fitness: math.Inf(1),
	}
}

func (ind Individual) Mutate() Individual {
	return NewIndividual(ind.Options)
}

func (ind Individual) Run(input []float64) (output []float64) {
	return []float64{0, 0, 0}
}
