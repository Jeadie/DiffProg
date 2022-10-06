package tensor

import "fmt"

type SymbolFunc interface {
	F(x []float64) []float64
	Df(x []float64) [][]float64
	ShapeIn() uint
	ShapeOut() uint
}

type tensorDiffFunc struct {
	diffFuncs []ScalarDiffFunc
	shapeIn   uint
}

func (d tensorDiffFunc) F(x []float64) []float64 {
	y := make([]float64, len(d.diffFuncs))
	for i, yi := range d.diffFuncs {
		y[i] = yi.F(x)[0]
	}
	return y
}
func (d tensorDiffFunc) Df(x []float64) [][]float64 {
	y := make([][]float64, len(d.diffFuncs))
	for i, yi := range d.diffFuncs {
		y[i] = yi.Df(x)[0]
	}
	return y
}
func (d tensorDiffFunc) ShapeIn() uint  { return d.shapeIn }
func (d tensorDiffFunc) ShapeOut() uint { return uint(len(d.diffFuncs)) }

func TensorDiffFunc(fnDims ...ScalarDiffFunc) SymbolFunc {
	return tensorDiffFunc{
		diffFuncs: fnDims,
		shapeIn:   uint(len(fnDims)),
	}
}

type constant struct {
	c []float64
}

func (c constant) F(x []float64) []float64 { return c.c }
func (c constant) Df(x []float64) [][]float64 {
	// zero matrix.
	y := make([][]float64, c.ShapeOut())
	for i := uint(0); i < c.ShapeOut(); i++ {
		y[i] = make([]float64, c.ShapeIn())
	}
	return y
}

func (c constant) ShapeIn() uint  { return uint(len(c.c)) }
func (c constant) ShapeOut() uint { return uint(len(c.c)) }

// Linear creates a multi-variate linear function.
func Linear(coefficients ...[]float64) SymbolFunc {
	if len(coefficients) <= 0 {
		return nil
	} else if len(coefficients[0]) == 1 {
		cs := make([]float64, len(coefficients))
		for i, c := range coefficients {
			cs[i] = c[0]
		}
		return constant{c: cs}

	} else if len(coefficients[0]) == 0 {
		panic("TODO")
	}

	return linear{
		coefficients: coefficients,
		shapeIn:      len(coefficients[0]) - 1,
		shapeOut:     len(coefficients),
	}
}

type linear struct {
	// coefficients ordered:
	// [y_0]  [ [a0, a1, a2, a3], [1]
	// [y_1]    [b0, b1, b2, b3], [x0]
	// [y_2]  	[c0, c1, c2, c3], [x1]
	// [y_3]  	[d0, d1, d2, d3]],[x2]
	coefficients [][]float64
	shapeIn      int
	shapeOut     int
}

func (p linear) ShapeIn() uint {
	return uint(p.shapeIn)
}

func (p linear) ShapeOut() uint {
	return uint(p.shapeOut)
}

// Df of a linear model is non-constant coefficients
func (p linear) Df(x []float64) [][]float64 {
	y := make([][]float64, len(p.coefficients))

	// Remove constants from derivative coefficients
	for i, c := range p.coefficients {
		yy := make([]float64, len(c)-1)
		copy(yy, c[1:])
		y[i] = yy
	}
	return y
}

// F applies a linear function to input x vector.
func (p linear) F(x []float64) []float64 {
	if len(x) != p.shapeIn {
		panic(fmt.Errorf("function has input shape %d. x.shape==%d", p.shapeIn, len(x)))
	}
	y := make([]float64, p.shapeOut)

	// Matrix multiply, coefficients * x
	for i, cDims := range p.coefficients {
		y[i] = dot(cDims[1:], x) + cDims[0]
	}
	return y
}

func dot(a []float64, b []float64) float64 {
	y := 0.0
	for i, _ := range a {
		y += a[i] * b[i]
	}
	return y
}

func scalarOp[T any](a []T, c T, op func(c, ai T) T) []T {
	y := make([]T, len(a))
	for i, _ := range a {
		y[i] = op(c, a[i])
	}
	return y
}

func scalarAdd(a []float64, b float64) []float64 {
	return scalarOp[float64](a, b, func(x float64, y float64) float64 { return x + y })
}

func scalarMul(a []float64, b float64) []float64 {
	return scalarOp[float64](a, b, func(x float64, y float64) float64 { return x * y })
}

func pairwiseOp[T any](a []T, b []T, op func(x T, y T) T) []T {
	if len(a) != len(b) {
		panic(fmt.Sprintf("vectors cannot be different lengths. %d != %d", len(a), len(b)))
	}
	y := make([]T, len(a))
	for i, _ := range a {
		y[i] = op(a[i], b[i])
	}
	return y
}

func vectorAdd(a []float64, b []float64) []float64 {
	return pairwiseOp[float64](a, b, func(x float64, y float64) float64 { return x + y })
}

func vectorMul(a []float64, b []float64) []float64 {
	return pairwiseOp[float64](a, b, func(x float64, y float64) float64 { return x * y })
}

func matrixAdd(a [][]float64, b [][]float64) [][]float64 {
	return pairwiseOp[[]float64](a, b, func(x []float64, y []float64) []float64 { return vectorAdd(x, y) })
}

// type TensorsDiffFunc struct {
// 	f        func(x []float64) []float64
// 	df       func(x []float64) [][]float64
// 	shapeOut uint
// }

// func Add(a TensorsDiffFunc, b TensorsDiffFunc) TensorsDiffFunc {
// 	return TensorsDiffFunc{
// 		f:  func(x []float64) []float64 { return vectorAdd(a.f(x), b.f(x)) },
// 		df: func(x []float64) [][]float64 { return matrixAdd(a.df(x), b.df(x)) },
// 	}
// }
