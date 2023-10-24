package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"time"
)

var randomRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = randomRunes[r.Intn(len(randomRunes))]
	}

	return string(b)
}

func GenerateSaltedHash(str string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	io.WriteString(hash, salt)
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateHashFromString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
