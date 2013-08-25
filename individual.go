package cgp

type Individual struct {
	Fitness float64
}

func NewIndividual(cgp *CGP) Individual {
	return Individual{}
}

func (ind Individual) Run(input []float64) (output []float64) {
	return
}
