package util

import (
	"crypto/rand"
	"math/big"
)

const LETTERS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func RandomString(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
		if err != nil {
			return "", err
		}
		ret[i] = LETTERS[num.Int64()]
	}

	return string(ret), nil
}
