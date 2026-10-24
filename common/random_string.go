package common

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString returns a cryptographically secure lowercase hexadecimal string.
func RandomString(n int) string {
	b := make([]byte, n/2+n%2)
	if _, err := rand.Read(b); err != nil {
		panic("secure random number generation failed: " + err.Error())
	}
	return hex.EncodeToString(b)[:n]
}
