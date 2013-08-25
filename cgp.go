package cgp

type CGPFunction func([]float64) float64
type EvalFunction func(Individual) float64
type RndConstFunction func() float64

type CGP struct {
}

func New(popSize uint, maxGenes uint, mutationRate float64, numInputs uint, numOutputs uint, maxArity uint, functionList []CGPFunction, randomConstant RndConstFunction, evaluator EvalFunction) *CGP {
	return &CGP{}
}
