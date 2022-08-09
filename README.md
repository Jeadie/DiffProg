# DiffProg
Minimal, differential programming in Go.  

## Symbolic Differentiation

Simple polynomials
```go
import (
    "fmt"
    f "github.com/Jeadie/DiffProg/pkg/scalar"
)

// f(x) = 1+3x-2x^2
poly := f.Polynomial(1.0, 3.0, 2.0)

fmt.Printf("f(0) = %f\n", poly.F(0.0)) // f(0) = 1.0
fmt.Printf("f'(0) = %f\n", poly.Df(0.0)) // f'(0) = 3.0
``` 

Common differentiable function with natural extension 
```go
// f(x) = sin(x) + cos(x)
f.Add(f.Sin, f.Cos)

// f(x) = sin(x) / e^x
f.Div(f.Sin, f.Exp)

// f(x) = sin(e^x) + e^cos(x)
fn := f.Add(
	f.Compose(f.Sin, f.Exp),
	f.Compose(f.Exp, f.Cos),
)

// f'(e^0.0)
fn.Df(math.Exp(0.0))
```

Run across vectors
```go
x := f.Linspace(-10.0, 10.0, 10)
y := f.Apply(fn.F, x)
dy := f.Apply(fn.Df, x)

fmt.Println(x)
fmt.Println(y)
fmt.Println(dy)
```
```
[-10 -8 -6 -4 -2 0 2 4 6 8]
[0.432 0.865 2.615 0.538 0.794 3.559 1.553 -0.408 3.577 1.266]
[-0.235 0.856 -0.727 -0.375 0.734 0.540 2.713 -19.844 106.727 -2730.648]
```

