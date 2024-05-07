package lib

import (
	"crypto/rand"
	"math/big"
)

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	letterBytes := []byte(letters)
	letterLen := big.NewInt(int64(len(letterBytes)))

	for i := 0; i < n; i++ {
		randIndex, _ := rand.Int(rand.Reader, letterLen)
		result += string(letterBytes[randIndex.Int64()])
	}
	return result
}
