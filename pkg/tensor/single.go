package tensor

import "math"

// ScalarDiffFunc represents a multi-variate differential equation with scalar output
// i.e. y = f(x1, x2, ...xn); y \in \R
type ScalarDiffFunc struct {
	f       func(x []float64) float64
	df      func(x []float64) []float64
	shapeIn uint
}

func (d ScalarDiffFunc) F(x []float64) []float64    { return []float64{d.f(x)} }
func (d ScalarDiffFunc) Df(x []float64) [][]float64 { return [][]float64{d.df(x)} }
func (d ScalarDiffFunc) ShapeIn() uint              { return d.shapeIn }
func (d ScalarDiffFunc) ShapeOut() uint             { return 1 }

// Sin function on k-th dimension of \R^n dimensional input.
func Sin(n, k uint) ScalarDiffFunc {
	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return math.Sin(x[k])
		},
		df: func(x []float64) []float64 {
			dy := make([]float64, n)
			dy[k] = math.Cos(x[k])
			return dy
		},
		shapeIn: n,
	}
}

// Cos function on k-th dimension of \R^n dimensional input.
func Cos(n, k uint) ScalarDiffFunc {
	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return math.Cos(x[k])
		},
		df: func(x []float64) []float64 {
			dy := make([]float64, n)
			dy[k] = -1.0 * math.Sin(x[k])
			return dy
		},
		shapeIn: n,
	}
}

func Lin(c float64, coefficients ...float64) ScalarDiffFunc {

	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return dot(coefficients, x) + c
		},
		df: func(x []float64) []float64 {
			return coefficients
		},
		shapeIn: uint(len(coefficients) - 1),
	}
}

// Exponential function on k-th dimension of \R^n dimensional input.
func Exp(n, k uint) ScalarDiffFunc {
	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return math.Exp(x[k])
		},
		df: func(x []float64) []float64 {
			dy := make([]float64, n)
			dy[k] = math.Exp(x[k])
			return dy
		},
		shapeIn: n,
	}
}

// Log function on k-th dimension of \R^n dimensional input.
func Log(n, k uint) ScalarDiffFunc {
	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return math.Log(x[k])
		},
		df: func(x []float64) []float64 {
			dy := make([]float64, n)
			dy[k] = 1.0 / x[k]
			return dy
		},
		shapeIn: n,
	}
}

func Add(a, b ScalarDiffFunc) ScalarDiffFunc {
	if a.ShapeIn() != b.ShapeIn() {
		panic("")
	}
	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return a.f(x) + b.f(x)
		},
		df: func(x []float64) []float64 {
			return vectorAdd(a.df(x), b.df(x))
		},
		shapeIn: a.ShapeIn(),
	}
}

// Mul two ScalarDiffFuncs
func Mul(a, b ScalarDiffFunc) ScalarDiffFunc {
	if a.ShapeIn() != b.ShapeIn() {
		panic("")
	}
	return ScalarDiffFunc{
		f: func(x []float64) float64 {
			return a.f(x) + b.f(x)
		},
		df: func(x []float64) []float64 {
			return vectorAdd(scalarMul(a.df(x), b.f(x)), scalarMul(b.df(x), a.f(x)))
		},
		shapeIn: a.ShapeIn(),
	}
}
