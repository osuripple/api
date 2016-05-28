package common

// In picks x if y < x, picks z if y > z, or if none of the previous
// conditions is satisfies, it simply picks y.
func In(x, y, z int) {
	switch {
	case y < x:
		return x
	case y > z:
		return z
	}
	return y
}
