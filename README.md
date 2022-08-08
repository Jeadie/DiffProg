# DiffProg
Differential Programming in Go 

## Symbolic Differentiation

Simple polynomials
```go
import (
    "fmt"
    f "github.com/Jeadie/DiffProg/pkg/diff"
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
	f.Compose(f.Exp, f.Cos)
)

// f'(e^0.0)
fn.Df(math.Exp(0.0))
```

Run across vectors
```go
x := f.Linspace(-10.0, 10.0, 100)
y := f.Apply(fn.Df, x)
```

