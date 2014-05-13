package neuralnetwork

import (
	"fmt"
	"math"
	"math/rand"
)

func random(a, b float64) float64 {
	return (b-a)*rand.Float64() + a
}

func matrix(I, J int) [][]float64 {
	m := make([][]float64, I)
	for i := 0; i < I; i++ {
		m[i] = make([]float64, J)
	}
	return m
}

func vector(I int, fill float64) []float64 {
	v := make([]float64, I)
	for i := 0; i < I; i++ {
		v[i] = fill
	}
	return v
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func dsigmoid(y float64) float64 {
	return y * (1 - y)
}

type NeuralNetwork struct {
	// Number of input, hidden and output nodes
	ni, nh, no int
	// Whether it is regression or not
	regression bool
	// Activations for nodes
	ai, ah, ao []float64
	// Weights
	wi, wo [][]float64
	// Last change in weights for momentum
	ci, co [][]float64
}

func New(ni, nh, no int, regression bool) *NeuralNetwork {
	nn := &NeuralNetwork{ni: ni + 1, nh: nh + 1, no: no, regression: regression}

	nn.ai = vector(nn.ni, 1.0)
	nn.ah = vector(nn.nh, 1.0)
	nn.ao = vector(nn.no, 1.0)

	nn.wi = matrix(nn.ni, nn.nh)
	nn.wo = matrix(nn.nh, nn.no)

	for i := 0; i < nn.ni; i++ {
		for j := 0; j < nn.nh; j++ {
			nn.wi[i][j] = random(-1, 1)
		}
	}

	for i := 0; i < nn.nh; i++ {
		for j := 0; j < nn.no; j++ {
			nn.wo[i][j] = random(-1, 1)
		}
	}

	nn.ci = matrix(nn.ni, nn.nh)
	nn.co = matrix(nn.nh, nn.no)

	return nn
}

func (nn *NeuralNetwork) Update(inputs []float64) []float64 {
	if len(inputs) != nn.ni-1 {
		fmt.Println("Error: wrong number of inputs")
		return []float64{} // should return error
	}

	for i := 0; i < nn.ni-1; i++ {
		nn.ai[i] = inputs[i]
	}

	for i := 0; i < nn.nh-1; i++ {
		var sum float64 = 0.0
		for j := 0; j < nn.ni; j++ {
			sum += nn.ai[j] * nn.wi[j][i]
		}
		nn.ah[i] = sigmoid(sum)
	}

	for i := 0; i < nn.no; i++ {
		var sum float64 = 0.0
		for j := 0; j < nn.nh; j++ {
			sum += nn.ah[j] * nn.wo[j][i]
		}
		if nn.regression {
			nn.ao[i] = sum
		} else {
			nn.ao[i] = sigmoid(sum)
		}
	}

	return nn.ao
}

func (nn *NeuralNetwork) BackPropagate(targets []float64, lRate, mFactor float64) float64 {
	if len(targets) != nn.no {
		fmt.Println("Error: wrong number of target values")
		return 0.0
	}

	output_deltas := vector(nn.no, 0.0)
	for i := 0; i < nn.no; i++ {
		output_deltas[i] = targets[i] - nn.ao[i]

		if !nn.regression {
			output_deltas[i] = dsigmoid(nn.ao[i]) * output_deltas[i]
		}
	}

	hidden_deltas := vector(nn.nh, 0.0)
	for i := 0; i < nn.nh; i++ {
		var e float64 = 0.0

		for j := 0; j < nn.no; j++ {
			e += output_deltas[j] * nn.wo[i][j]
		}
		hidden_deltas[i] = dsigmoid(nn.ah[i]) * e
	}

	for i := 0; i < nn.nh; i++ {
		for j := 0; j < nn.no; j++ {
			change := output_deltas[j] * nn.ah[i]
			nn.wo[i][j] = nn.wo[i][j] + lRate*change + mFactor*nn.co[i][j]
			nn.co[i][j] = change
		}
	}

	for i := 0; i < nn.ni; i++ {
		for j := 0; j < nn.nh; j++ {
			change := hidden_deltas[j] * nn.ai[i]
			nn.wi[i][j] = nn.wi[i][j] + lRate*change + mFactor*nn.ci[i][j]
			nn.ci[i][j] = change
		}
	}

	var e float64 = 0.0

	for i := 0; i < len(targets); i++ {
		e += 0.5 * math.Pow(targets[i]-nn.ao[i], 2)
	}

	return e
}

func (nn *NeuralNetwork) Train(patterns [][][]float64, iterations int, lRate, mFactor float64) []float64 {
	errors := make([]float64, iterations)

	for i := 0; i < iterations; i++ {
		var e float64 = 0.0
		for _, p := range patterns {
			nn.Update(p[0])

			tmp := nn.BackPropagate(p[1], lRate, mFactor)
			e += tmp
		}

		errors[i] = e
	}

	return errors
}

func (nn *NeuralNetwork) Test(patterns [][][]float64) {
	for _, p := range patterns {
		fmt.Println(p[0], "->", nn.Update(p[0]), " : ", p[1])
	}
}
