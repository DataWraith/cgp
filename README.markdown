
# CGP for Go

This is a library implementing [Cartesian Genetic Programming](http://www.cartesiangp.co.uk/)
in the Go programming language.

## Usage

Here's an example of using CGP for symbolic regression. We are trying to
approximate the function f(x) = x³ - 2x + 10.

```go
    import "github.com/DataWraith/cgp"
    import "math"
    import "math/rand"

    // First, we set up our parameters:
    options := cgp.CGPOptions{
            PopSize:      100,  // The population size. One parent, 99 children
            NumGenes:     30,   // The maximum number of functions in the genome
            MutationRate: 0.05, // The mutation rate
            NumInputs:    1,    // The number of input values
            NumOutputs:   1,    // The number of output values
            MaxArity:     2,    // The maximum arity of the functions in the FunctionList
    }

    // The function list specifies the functions that are used in the genetic
    // program. The first input is always the constant evolved for the function.
    options.FunctionList = []cgp.CGPFunction{
            func(input []float64) float64 { return input[0] },            // Constant
            func(input []float64) float64 { return input[1] },            // Pass through A
            func(input []float64) float64 { return input[2] },            // Pass through B
            func(input []float64) float64 { return input[1] + input[2] }, // Addition
            func(input []float64) float64 { return input[1] - input[2] }, // Subtraction
            func(input []float64) float64 { return input[1] * input[2] }, // Multiplication

            // Division
            func(input []float64) float64 {
                    if input[2] == 0 {
                            return 0
                    }
                    return input[1] / input[2]
            },
    }

    // This function is used to generate constants for use by the functions in
    // the FunctionList. In this case we are using the provided instance of
    // rand.Rand to generate a random constant between 0 and 100.
    options.RandConst = func(rand *rand.Rand) float64 {
            return float64(rand.Intn(101))
    }

    // Prepare the testcases.
    var testCases = []struct {
            in  float64
            out float64
    }{
            {0, 10},
            {0.5, 9.125},
            {1, 9},
            {10, 990},
            {-5, -105},
            {17, 4889},
            {3.14, 34.679144},
    }

    // The evaluator uses the test cases to grade an individual by setting the
    // fitness to the sum of squared errors. The lower the fitness the better the
    // individual. Note that the input to the individual is a slice of float64.
    options.Evaluator = func(ind cgp.Individual) float64 {
           fitness := 0.0
           for _, tc := range testCases {
                    input := []float64{tc.in}
                    output := ind.Run(input)
                    fitness += math.Pow(tc.out-output[0], 2)
            }
            return fitness
    }

    // Initialize CGP
    gp := cgp.New(options)

    // Population[0] is the parent, which is the most fit individual. We
    // loop until we've found a perfect solution (fitness 0)
    for gp.Population[0].Fitness > 0 {
            gp.RunGeneration()
            fmt.Println(gp.Population[0].Fitness)
    }
```
