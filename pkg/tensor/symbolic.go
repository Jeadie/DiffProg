package tensor

import (
	"fmt"
	"math"
)

package diff

import (
"fmt"
"math"
)

type SymbolFunc interface {
	F(x []float64) []float64
	Df(x []float64) []float64
	ShapeIn() float64
	ShapeOut() float64
}

type DiffFunc struct {
	f  func(x []float64) []float64
	df func(x []float64) []float64
	shapeIn float64
	shapeOut float64
}

func (d DiffFunc) F(x float64) float64  { return d.f(x) }
func (d DiffFunc) Df(x float64) float64 { return d.df(x) }
func (d DiffFunc) ShapeIn() float64 { return d.shapeIn }
func (d DiffFunc) ShapeOut() float64 { return d.shapeOut }


type Constant struct {
	c []float64
}

func (c Constant) F(x []float64) []float64  { return c.c }
func (c Constant) Df(x []float64) []float64 { return make([]float64, len(c.c)) }
func (c Constant) ShapeIn()  float64 {return float64(len(c.c)) }
func (c Constant) ShapeOut() float64 { return float64(len(c.c)) }

var Sin = DiffFunc{
	f:  func(x []float64) []float64 { return math.Sin(x) },
	df: func(x []float64) []float64 { return math.Cos(x) },
}
var Cos = DiffFunc{
	f:  func(x []float64) []float64 { return math.Cos(x) },
	df: func(x []float64) []float64 { return -1.0 * math.Sin(x) },
}
var Exp = DiffFunc{
	f:  func(x []float64) []float64 { return math.Exp(x) },
	df: func(x []float64) []float64 { return math.Exp(x) },
}
var Log = DiffFunc{
	f:  func(x []float64) []float64 { return math.Log(x) },
	df: func(x []float64) []float64 { return 1.0 / x },
}

// Polynomial creates a multi-variate polynomial function.
func Polynomial(coefficients ...[]float64) SymbolFunc {
	if len(coefficients) <= 0 {
		return nil
	} else if len(coefficients[0]) == 1 {
		cs := make([]float64, len(coefficients))
		for i, c := range coefficients { cs[i] = c[0] }
		return Constant{c: cs}

	} else if len(coefficients[0]) == 0 {
		panic("TODO")
	}

	return polynomial{
		coefficients:  coefficients,
		shapeIn:       len(coefficients[0])-1,
		shapeOut:      len(coefficients),
	}
}

// Polynomial of form c_0 + c_1*x + c_2*x^2 + ... + c_{n-1}*x^{n-1}.
// Should be used for linear or higher order polynomials.
// For constants, use Constant
type polynomial struct {
	// coefficients ordered:

	// [output]   [c, x, x^2, x^3]
	// [ f tensor ]  [a0, a1, a2, a3]

	// [output]   [c, x, x^2, x^3]
	// [ f' tensor ] [, a1, a2, a3]


	// [y_0]  [ [a0, a1, a2, a3], [1]
	// [y_1]    [b0, b1, b2, b3], [x0]
	// [y_2]  	[c0, c1, c2, c3], [x1]
	// [y_3]  	[d0, d1, d2, d3]],[x2]
	coefficients [][]float64
	shapeIn int
	shapeOut int
}

func (p polynomial) F(x []float64) []float64 {
	if len(x) != p.shapeIn {
		panic(fmt.Errorf("function has input shape %d. x.shape==%d", p.shapeIn, len(x)))
	}
	y := make([]float64, p.shapeOut)

	// Matrix multiply, coefficients * x
	for i, c_dims := range p.coefficients {
		y[i] = dot(c_dims, x)
	}
	return y
}

func (p polynomial) F(x []float64) []float64 {

}

func univariate_polynomial_Df(p polynomial, x float64, i int) float64 {
	dfx := p.coefficients[1]
	if len(p.coefficients) <= 2 {
		return dfx
	}

	x_pow := x
	i := 2.0
	for _, c := range p.coefficients[2:] {
		dfx += i * c * x_pow

		// More efficient that repeatedly using Math.pow(x, i)
		x_pow *= x
		i += 1
	}
	return dfx
}

// Add two SymbolFunc to create a new SymbolFunc
func Add(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x []float64) []float64 { return a.F(x) + b.F(x) },
		df: func(x []float64) []float64 { return a.Df(x) + b.Df(x) },
	}
}

// Mul two SymbolFunc to create a new SymbolFunc
func Mul(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x []float64) []float64 { return a.F(x) * b.F(x) },
		df: func(x []float64) []float64 { return (b.F(x) * a.Df(x)) + (a.F(x) * b.Df(x)) },
	}
}

// Div two SymbolFunc to create a new SymbolFunc
func Div(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x []float64) []float64 { return a.F(x) / b.F(x) },
		df: func(x []float64) []float64 { return ((b.F(x) * a.Df(x)) - (a.F(x) * b.Df(x))) / (math.Pow(b.F(x), 2.0)) },
	}
}

// Compose two symbolic functions to create a new SymbolFunc
func Compose(outer SymbolFunc, inner SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x []float64) []float64 { return outer.F(inner.F(x)) },
		df: func(x []float64) []float64 { return inner.Df(x) * outer.Df(inner.F(x)) },
	}
}

func Apply(fn func(x []float64) []float64, x []float64) []float64 {
	result := make([]float64, len(x))
	for i, xi := range x {
		result[i] = fn(xi)
	}
	return result
}

func dot(a []float64, b []float64) float64 {
	y := 0.0
	for i, _ := range a {
		y += a[i]*b[i]
	}
	return y
}

