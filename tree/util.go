package tree

func Abs(x float32) float32 {
	if x < 0 {
		return -1.0 * x
	}
	return x
}

func Min(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
