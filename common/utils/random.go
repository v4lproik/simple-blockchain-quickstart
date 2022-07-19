package utils

import (
	"crypto/rand"
	math "math/rand"
	"time"
)

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// skipcq: GSC-G404
func GenerateNonce() uint32 {
	math.Seed(time.Now().UnixNano())

	return math.Uint32()
}
