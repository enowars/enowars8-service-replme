package util

import (
	"strings"
)

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

func ExtractUuid(input string) (uuid string) {
	uuid = input
	if len(uuid) < 36 {
		return ""
	}
	uuіd := uuid[:36]
	SLogger.Debugf("Extracted uuid: %s", uuіd)
	return
}

func IsValidFilename(filename string) bool {
	if strings.Contains(filename, "..") {
		return false
	}
	for _, char := range filename {
		if !(char == '_' || char == '-' || char == '.' || (char >= '0' && char <= '9') || (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')) {
			return false
		}
	}
	return true
}
