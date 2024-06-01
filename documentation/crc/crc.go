package main

import (
	"crypto/rand"
	"fmt"
	"hash/crc32"
	"io"
	"math/big"
	"net/http"
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

func BruteForceCRC32(input string) string {
	max := 16

	u := []byte(input)
	goal := crc32.ChecksumIEEE(u)
	fmt.Printf("[Goal] : %s (CRC: %d)\n", string(u), goal)

	result := ""

	for i := 1; i <= max; i++ {
		if result != "" {
			break
		}
		fmt.Printf("[Step] : Length: %d ...\n", i)
		var wg sync.WaitGroup
		for j :=0; j <= 8; j++ {
			wg.Add(1)
			go func(i int, num int, step int) {
				defer wg.Done()
				buf := make([]byte, i)
				for j := 0; j < i; j++ {
					buf[j] = 'a'
				}
				buf[0] += byte(num)

				for {
					test := crc32.ChecksumIEEE(buf)
					if test == goal && string(buf) != input {
						fmt.Println("[MATCH]:", string(buf))
						result = string(buf)
						return
					}

					if result != "" {
						return
					}

					var j int
					for j = 0; j < i; j++ {
						if j == 0 && buf[j] <= ('z'-byte(step)) {
							buf[j] += byte(step)
							break
						}
						if j == 0 && buf[j] > ('z'-byte(step)) {
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
					if j >= i {
						break
					}
				}
			}(i, j, 8)
		}
		wg.Wait()
	}
	return result
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	input, _ := io.ReadAll(r.Body)
	result := BruteForceCRC32(string(input))
	io.WriteString(w, result)
}

func main() {

	http.HandleFunc("/", HandleRequest)
	http.ListenAndServe(":3333", nil)
}
