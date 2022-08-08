package diff

func Linspace(start float64, stop float64, num uint) []float64 {
	if stop < start {
		return []float64{}
	}

	result := make([]float64, num)
	x := start
	dx := (stop - start) / float64(num)

	for i := uint(0); i < num; i++ {
		result[i] = x
		x += dx
	}
	return result
}
