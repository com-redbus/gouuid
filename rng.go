package gouuid

import (
	"crypto/rand"
)

func rng(n int) ([]byte, error) {
	return generateRandomBytes(n)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
