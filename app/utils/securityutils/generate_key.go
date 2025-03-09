package securityutils

import (
	"math/rand"
	"time"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
)

func GenerateKey(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = letters[rand.Intn(len(letters))]
		} else {
			b[i] = digits[rand.Intn(len(digits))]
		}
	}
	return string(b)
}
