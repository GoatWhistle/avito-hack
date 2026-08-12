package domain

import (
	"crypto/rand"
	"strings"
)

const (
	displayIDLength   = 12
	displayIDByteSize = 8
	crockfordAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"
	displayIDAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func NewDisplayID() string {
	buf := make([]byte, displayIDByteSize)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	return encodeCrockfordBase32(buf)
}

func IsValidDisplayID(value string) bool {
	if len(value) != displayIDLength {
		return false
	}

	for _, r := range value {
		if !strings.ContainsRune(displayIDAlphabet, r) {
			return false
		}
	}

	return true
}

func encodeCrockfordBase32(data []byte) string {
	var sb strings.Builder
	sb.Grow(displayIDLength)

	var bitBuffer uint64
	bitCount := 0

	for _, b := range data {
		bitBuffer = (bitBuffer << 8) | uint64(b)
		bitCount += 8

		for bitCount >= 5 {
			bitCount -= 5
			index := (bitBuffer >> uint(bitCount)) & 0x1F
			sb.WriteByte(crockfordAlphabet[index])
		}
	}

	if bitCount > 0 {
		index := (bitBuffer << uint(5-bitCount)) & 0x1F
		sb.WriteByte(crockfordAlphabet[index])
	}

	encoded := sb.String()
	if len(encoded) < displayIDLength {
		encoded += strings.Repeat("0", displayIDLength-len(encoded))
	}

	return encoded[:displayIDLength]
}
