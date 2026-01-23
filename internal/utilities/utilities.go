package utilities

import "math/rand/v2"

// Roll calculates a random number in the range of lower, upper parameters (inclusive)
func Roll(lower int, upper int) int {
	if upper < lower {
		panic("Roll: upper < lower")
	}
	return lower + rand.IntN(1+upper-lower)
}
