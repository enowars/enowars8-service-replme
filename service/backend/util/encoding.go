package util

func DecodeSpecialChars(input []byte) []byte {
	ret := make([]byte, len(input))
	for i, c := range input {
		if c == 0b110000 {
			ret[i] = 0b1100000
		} else if c > 0b110000 && c < 0b110110 {
			ret[i] = c + 74
		} else {
			ret[i] = c
		}
	}
	return ret
}

