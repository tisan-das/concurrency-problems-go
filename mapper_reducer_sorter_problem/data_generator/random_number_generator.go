package data_generator

import (
	"math/rand"
)

type RandomGenerator struct {
	seedValue int
}

func NewRandomNumberGenerator(seed int) *RandomGenerator {
	if seed < 0 {
		seed = rand.Int()
	}
	rand.Seed(int64(seed))
	return &RandomGenerator{
		seedValue: seed,
	}
}

func (generator *RandomGenerator) GenerateRandomInt() int {
	return rand.Int()
}
