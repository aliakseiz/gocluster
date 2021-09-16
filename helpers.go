package cluster

import "math"

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

// count number of digits, for example 123356 will return 6.
func digitsCount(a int) int {
	return int(math.Floor(math.Log10(math.Abs(float64(a))))) + 1
}

var epsilon = 0.0000000001

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
}
