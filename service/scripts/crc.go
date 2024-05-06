package main

import (
	"fmt"
	"crypto/rand"
	"math/big"
	"hash/crc32"
	"sync"
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

func main() {

	max := 16

	str := RandomString(32)
	// str := "test"
	u := []byte(str)
	goal := crc32.ChecksumIEEE(u)
	fmt.Printf("[Goal] : %s (CRC: %d)\n", string(u), goal)

	for i := 1; i <= max; i++ {
		fmt.Printf("[Step] : Length: %d ...\n", i)
		var wg sync.WaitGroup
		for j := range 8 {
			wg.Add(1)
			go func (i int, num int, step int) {
				defer wg.Done()
				buf := make([]byte, i)
				for j := 0; j < i; j++ {
					buf[j] = 'a'
				}
				buf[0] += byte(num)

				for {
					test := crc32.ChecksumIEEE(buf)
					if test == goal {
						fmt.Println("[MATCH]:", string(buf))
					}

					var j int
					for j = 0; j < i; j++ {
						if j == 0 && buf[j] <= ('z' - byte(step)) {
							buf[j] += byte(step)
							break
						}
						if j == 0 && buf[j] > ('z' - byte(step)) {
							buf[j] = 'a' + byte(num)
							continue
						}
						if buf[j] < 'z' {
							buf[j]++
							break
						}
						if buf[j] >= 'z' {
							buf[j] = 'a'
						}
					}
					if (j >= i) {
						break;
					}
				}
			}(i, j, 8)
		}
		wg.Wait()
	}
}

