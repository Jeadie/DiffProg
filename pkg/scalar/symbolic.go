package scalar

import "math"

type SymbolFunc interface {
	F(x float64) float64
	Df(x float64) float64
}

type DiffFunc struct {
	f  func(x float64) float64
	df func(x float64) float64
}

func (d DiffFunc) F(x float64) float64  { return d.f(x) }
func (d DiffFunc) Df(x float64) float64 { return d.df(x) }

type Constant struct {
	c float64
}

func (c Constant) F(x float64) float64  { return c.c }
func (c Constant) Df(x float64) float64 { return 0.0 }

var Sin = DiffFunc{
	f:  func(x float64) float64 { return math.Sin(x) },
	df: func(x float64) float64 { return math.Cos(x) },
}
var Cos = DiffFunc{
	f:  func(x float64) float64 { return math.Cos(x) },
	df: func(x float64) float64 { return -1.0 * math.Sin(x) },
}
var Exp = DiffFunc{
	f:  func(x float64) float64 { return math.Exp(x) },
	df: func(x float64) float64 { return math.Exp(x) },
}
var Log = DiffFunc{
	f:  func(x float64) float64 { return math.Log(x) },
	df: func(x float64) float64 { return 1.0 / x },
}

func Polynomial(coefficients ...float64) SymbolFunc {
	return polynomial{coefficients: coefficients}
}

// Polynomial of form c_0 + c_1*x + c_2*x^2 + ... + c_{n-1}*x^{n-1}.
// Should be used for linear or higher order polynomials.
// For constants, use Constant
type polynomial struct {
	// coefficients ordered: c_0, c_1, c_2,..., c_{n-1}
	coefficients []float64
}

func (p polynomial) F(x float64) float64 {
	fx := p.coefficients[0]

	x_pow := x
	for _, c := range p.coefficients[1:] {
		fx += c * x_pow

		// More efficient that repeatedly using Math.pow(x, i)
		x_pow *= x
	}
	return fx
}

func (p polynomial) Df(x float64) float64 {
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
		f:  func(x float64) float64 { return a.F(x) + b.F(x) },
		df: func(x float64) float64 { return a.Df(x) + b.Df(x) },
	}
}

// Mul two SymbolFunc to create a new SymbolFunc
func Mul(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x float64) float64 { return a.F(x) * b.F(x) },
		df: func(x float64) float64 { return (b.F(x) * a.Df(x)) + (a.F(x) * b.Df(x)) },
	}
}

// Div two SymbolFunc to create a new SymbolFunc
func Div(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x float64) float64 { return a.F(x) / b.F(x) },
		df: func(x float64) float64 { return ((b.F(x) * a.Df(x)) - (a.F(x) * b.Df(x))) / (math.Pow(b.F(x), 2.0)) },
	}
}

// Compose two symbolic functions to create a new SymbolFunc
func Compose(outer SymbolFunc, inner SymbolFunc) SymbolFunc {
	return DiffFunc{
		f:  func(x float64) float64 { return outer.F(inner.F(x)) },
		df: func(x float64) float64 { return inner.Df(x) * outer.Df(inner.F(x)) },
	}
}

func Apply(fn func(x float64) float64, x []float64) []float64 {
	result := make([]float64, len(x))
	for i, xi := range x {
		result[i] = fn(xi)
	}
	return result
}
