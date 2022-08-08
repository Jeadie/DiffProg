package diff

import (
	"math"
	"testing"
)

type NumericTest struct {
	dfn    DFn
	dfnStr string
	x      []float64
	f      []float64
	df     []float64
}

var numericTests = []NumericTest{
	{
		dfn:    Constant{c: 5.0}.f,
		dfnStr: "y=5.0",
		x:      []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:      []float64{5.0, 5.0, 5.0, 5.0, 5.0},
		df:     []float64{0.0, 0.0, 0.0, 0.0, 0.0},
	},
	{
		dfn:    Polynomial{coefficients: []float64{4.0, 5.0}}.f,
		dfnStr: "y=5x+4",
		x:      []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:      []float64{-6.0, -1.0, 4.0, 9.0, 14.0},
		df:     []float64{5.0, 5.0, 5.0, 5.0, 5.0},
	},
	{
		dfn:    Polynomial{coefficients: []float64{4.0, 0.0, 1.0}}.f,
		dfnStr: "y=x^2+4",
		x:      []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:      []float64{8.0, 5.0, 4.0, 5.0, 8.0},
		df:     []float64{-4.0, -2.0, 0.0, 2.0, 4.0},
	},
	{
		dfn: Add(
			Polynomial{coefficients: []float64{0.0, 0.0, 1.0}},
			Constant{c: 4.0}).F,
		dfnStr: "y=x^2+4 (Add)",
		x:      []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:      []float64{8.0, 5.0, 4.0, 5.0, 8.0},
		df:     []float64{-4.0, -2.0, 0.0, 2.0, 4.0},
	},
	{
		dfn: Add(
			Mul(Polynomial{coefficients: []float64{0.0, 1.0}}, Polynomial{coefficients: []float64{0.0, 1.0}}),
			Constant{c: 4.0}).F,
		dfnStr: "y=x^2+4 (Mul + Add)",
		x:      []float64{-2.0, -1.0, 0.0, 1.0, 2.0},
		f:      []float64{8.0, 5.0, 4.0, 5.0, 8.0},
		df:     []float64{-4.0, -2.0, 0.0, 2.0, 4.0},
	},
}

func TestNumeric(t *testing.T) {
	for _, test := range numericTests {
		t.Run(test.dfnStr, func(t *testing.T) {
			fn := test.dfn
			centralDf := CentralDerive(fn, 0.0001)
			backwardDf := BackwardDerive(fn, 0.0001)
			ForwardDf := ForwardDerive(fn, 0.0001)

			for i, x := range test.x {
				central := centralDf(x)
				back := backwardDf(x)
				forward := ForwardDf(x)

				if math.Abs(test.df[i]-central) > 0.0002 {
					t.Errorf("Central numeric derivative expected %f != %f actual", test.df[i], central)
				}
				if math.Abs(test.df[i]-central) > 0.0002 {
					t.Errorf("Back numeric derivative expected %f != %f actual", test.df[i], back)
				}
				if math.Abs(test.df[i]-forward) > 0.0002 {
					t.Errorf("Forward numeric derivative expected %f != %f actual", test.df[i], forward)
				}
			}
		})
	}
}
