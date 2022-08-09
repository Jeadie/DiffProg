package scalar

type DFn func(x float64) float64

// ForwardDerive a function based on a given spacing, h and forward differencing.
func ForwardDerive(fn DFn, h float64) DFn {
	return func(x float64) float64 {
		return (fn(x+h) - fn(x)) / h
	}
}

// BackwardDerive a function based on a given spacing, h and backward differencing.
func BackwardDerive(fn DFn, h float64) DFn {
	return func(x float64) float64 {
		return (fn(x) - fn(x-h)) / h
	}
}

// CentralDerive a function based on a given spacing, h and central differencing.
func CentralDerive(fn DFn, h float64) DFn {
	step := h / 2
	return func(x float64) float64 {
		return (fn(x+step) - fn(x-step)) / (2 * step)
	}
}
