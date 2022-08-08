package diff

import "testing"

type SymbolicTest struct {
	symbol    SymbolFunc
	symbolStr string
	x         []float64
	f         []float64
	df        []float64
}

var tests = []SymbolicTest{
	{
		symbol:    Constant{c: 5.0},
		symbolStr: "y=5.0",
		x:         []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:         []float64{5.0, 5.0, 5.0, 5.0, 5.0},
		df:        []float64{0.0, 0.0, 0.0, 0.0, 0.0},
	},
	{
		symbol:    Polynomial{coefficients: []float64{4.0, 5.0}},
		symbolStr: "y=5x+4",
		x:         []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:         []float64{-6.0, -1.0, 4.0, 9.0, 14.0},
		df:        []float64{5.0, 5.0, 5.0, 5.0, 5.0},
	},
	{
		symbol:    Polynomial{coefficients: []float64{4.0, 0.0, 1.0}},
		symbolStr: "y=x^2+4",
		x:         []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:         []float64{8.0, 5.0, 4.0, 5.0, 8.0},
		df:        []float64{-4.0, -2.0, 0.0, 2.0, 4.0},
	},
	{
		symbol: Add(
			Polynomial{coefficients: []float64{0.0, 0.0, 1.0}},
			Constant{c: 4.0}),
		symbolStr: "y=x^2+4 (Add)",
		x:         []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:         []float64{8.0, 5.0, 4.0, 5.0, 8.0},
		df:        []float64{-4.0, -2.0, 0.0, 2.0, 4.0},
	},
	{
		symbol: Add(
			Mul(Polynomial{coefficients: []float64{0.0, 1.0}}, Polynomial{coefficients: []float64{0.0, 1.0}}),
			Constant{c: 4.0}),
		symbolStr: "y=x^2+4 (Mul + Add)",
		x:         []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:         []float64{8.0, 5.0, 4.0, 5.0, 8.0},
		df:        []float64{-4.0, -2.0, 0.0, 2.0, 4.0},
	},
}

func TestSymbolic(t *testing.T) {
	for _, test := range tests {
		t.Run(test.symbolStr, func(t *testing.T) {
			symbl := test.symbol

			for i, x := range test.x {
				symblF := symbl.F(x)
				if test.f[i] != symblF {
					t.Errorf("function expected %f != %f actual", test.f[i], symblF)
				}
				symblDf := symbl.df(x)
				if test.df[i] != symblDf {
					t.Errorf("derivative expected %f  != %f actual", test.df[i], symblDf)
				}
			}
		})
	}
}
