package main

import "math/rand"

// Vector2D is a struct that stores X and Y values for a position
type Vector2D struct {
	X float64
	Y float64
}

// Function that returns randomly either a or b. If a and b are equal, it returns value 'a'.
func randomChoice(a, b float64) float64 {
	if a == b {
		return a
	}
	if rand.Intn(2) == 0 {
		return a
	}
	return b
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
