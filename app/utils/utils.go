package utils

import "crypto/rand"

//gen salt of length n
func GenSalt(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
