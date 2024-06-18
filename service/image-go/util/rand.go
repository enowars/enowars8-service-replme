package util

import (
	"crypto/rand"
	"math/big"
	"os"
)

const LETTERS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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

func ApiKey(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		apikey, _ := RandomString(64)
		err := os.WriteFile(path, []byte(apikey), 0600)
		if err != nil {
			panic(err)
		}
		return apikey
	}

	return string(bytes)
}
