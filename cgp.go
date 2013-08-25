package cgp

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
	}

	for i := uint(0); i < result.PopSize; i++ {
		result.Population = append(result.Population, NewIndividual(result))
	}

	return result
}

func (cgp *CGP) RunGeneration() {
}
