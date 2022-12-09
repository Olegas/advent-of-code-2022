package util

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func AbsDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
