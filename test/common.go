package test

import (
	"crypto/rand"
	"encoding/hex"
	"math"
)

func RandomString(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}

func RandomEmail() string {
	return RandomString(8) + "@example.com"
}

func RandomPassword() string {
	return RandomString(16)
}
