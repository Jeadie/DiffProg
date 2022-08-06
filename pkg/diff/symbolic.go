package diff

import "math"

type SymbolFunc interface {
	f(x float64) float64
	df(x float64) float64
}

type DiffFunc struct {
	_f  func(x float64) float64
	_df func(x float64) float64
}

func (d DiffFunc) f(x float64) float64 {
	return d._f(x)
}

func (d DiffFunc) df(x float64) float64 {
	return d._df(x)
}

type Constant struct {
	c float64
}

var Sin = DiffFunc{
	_f:  func(x float64) float64 { return math.Sin(x) },
	_df: func(x float64) float64 { return math.Cos(x) },
}
var Cos = DiffFunc{
	_f:  func(x float64) float64 { return math.Cos(x) },
	_df: func(x float64) float64 { return -1.0 * math.Sin(x) },
}
var Exp = DiffFunc{
	_f:  func(x float64) float64 { return math.Exp(x) },
	_df: func(x float64) float64 { return math.Exp(x) },
}
var Log = DiffFunc{
	_f:  func(x float64) float64 { return math.Log(x) },
	_df: func(x float64) float64 { return 1.0 / x },
}

func (c Constant) f(x float64) float64 {
	return c.c
}

func (c Constant) df(x float64) float64 {
	return 0.0
}

// Polynomial of form c_0 + c_1*x + c_2*x^2 + ... + c_{n-1}*x^{n-1}.
// Should be used for linear or higher order polynomials.
// For constants, use Constant
type Polynomial struct {
	// coefficients ordered: c_0, c_1, c_2,..., c_{n-1}
	coefficients []float64
}

func (p Polynomial) f(x float64) float64 {
	fx := p.coefficients[0]

	x_pow := x
	for _, c := range p.coefficients[1:] {
		fx += c * x_pow

		// More efficient that repeatedly using Math.pow(x, i)
		x_pow *= x
	}
	return fx
}

func (p Polynomial) df(x float64) float64 {
	// of the form c_1 + 2*c_2*x + ... + (n-1)*c_{n-1}*x^{n-2}

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
		_f:  func(x float64) float64 { return a.f(x) + b.f(x) },
		_df: func(x float64) float64 { return a.df(x) + b.df(x) },
	}
}

// Mul two SymbolFunc to create a new SymbolFunc
func Mul(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		_f:  func(x float64) float64 { return a.f(x) * b.f(x) },
		_df: func(x float64) float64 { return (b.f(x) * a.df(x)) + (a.f(x) * b.df(x)) },
	}
}

// Div two SymbolFunc to create a new SymbolFunc
func Div(a SymbolFunc, b SymbolFunc) SymbolFunc {
	return DiffFunc{
		_f:  func(x float64) float64 { return a.f(x) / b.f(x) },
		_df: func(x float64) float64 { return ((b.f(x) * a.df(x)) - (a.f(x) * b.df(x))) / (math.Pow(b.f(x), 2.0)) },
	}
}

// Compose two symbolic functions to create a new SymbolFunc
func Compose(inner SymbolFunc, outer SymbolFunc) SymbolFunc {
	return DiffFunc{
		_f:  func(x float64) float64 { return outer.f(inner.f(x)) },
		_df: func(x float64) float64 { return inner.df(x) * outer.df(inner.f(x)) },
	}
}
