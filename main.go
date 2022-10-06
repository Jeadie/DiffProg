package main

import (
	"fmt"

	f "github.com/Jeadie/DiffProg/pkg/tensor"
)

func main() {
	// f(x, y) = {2x + 4y, 1 + y}
	fn := f.Linear([]float64{0.0, 2.0, 4.0}, []float64{1.0, 0.0, 1.0})
	y := fn.F([]float64{0.0, 0.0})
	dy := fn.Df([]float64{0.0, 0.0})
	fmt.Println(y)
	fmt.Println(dy)
	fmt.Println("---")
	// f(x, y) = 2x + 4y
	fn_1 := f.Lin(0.0, 2.0, 4.0)

	// f(x, y) = 1 + y
	fn_2 := f.Lin(1.0, 0.0, 1.0)

	// f(x, y) = {2x + 4y, 1 + y}
	fn = f.TensorDiffFunc(fn_1, fn_2)
	y = fn.F([]float64{0.0, 0.0})
	dy = fn.Df([]float64{0.0, 0.0})
	fmt.Println(y)
	fmt.Println(dy)

	f.Mul(fn_2, f.Cos(fn_2.ShapeIn(), 0))

}
