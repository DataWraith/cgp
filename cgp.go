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
	}

	for i := uint(0); i < result.PopSize; i++ {
		result.Population = append(result.Population, NewIndividual(result))
	}

	return result
}

func (cgp *CGP) RunGeneration() {
}
