package cluster

import "math"

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}

	return int(val + 0.5)
}

// digitsCount counts number of digits.
// 123456 will return 6.
func digitsCount(a int) int {
	return int(math.Floor(math.Log10(math.Abs(float64(a))))) + 1
}
