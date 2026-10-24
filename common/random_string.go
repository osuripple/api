package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
)

// LetterChars contains the lowercase and uppercase ASCII letters.
const LetterChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomString returns a cryptographically secure lowercase hexadecimal string.
func RandomString(n int) (string, error) {
	b := make([]byte, n/2+n%2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:n], nil
}

// RandomLetters returns a cryptographically secure string of n ASCII letters.
func RandomLetters(n int) (string, error) {
	return RandomStringFromChars(n, LetterChars)
}

// RandomStringFromChars returns a cryptographically secure string of n
// characters, each picked uniformly from chars. chars is treated as a sequence
// of bytes, so it should only contain ASCII characters.
func RandomStringFromChars(n int, chars string) (string, error) {
	if len(chars) == 0 {
		return "", errors.New("common: RandomStringFromChars called with empty chars")
	}
	max := big.NewInt(int64(len(chars)))
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = chars[idx.Int64()]
	}
	return string(b), nil
}
