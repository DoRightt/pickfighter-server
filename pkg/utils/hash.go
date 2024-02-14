package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"time"
)

// random string rune letter values
var randomRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GetRandomString returns random symbol string with specified length = n
func GetRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = randomRunes[r.Intn(len(randomRunes))]
	}

	return string(b)
}

// GenerateSaltedHash generates a SHA-256 hash with a salt for the given string.
// It returns the hash in hexadecimal format.
func GenerateSaltedHash(str string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	io.WriteString(hash, salt)
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateHashFromString generates a SHA-256 hash for the given string.
// It returns the hash in hexadecimal format.
func GenerateHashFromString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
