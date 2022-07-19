package utils

import (
	"crypto/rand"
	math "math/rand"
)

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateNonce() uint32 {
	math.Seed(DefaultTimeService.Nano())

	// skipcq
	return math.Uint32()
}
