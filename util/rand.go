package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandString simply generates random string of length n
func RandString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandEmail generates rand email
func RandEmail() string {
	return RandString(10) + "@" + RandString(5) + "." + RandString(3)
}

// RandString10 generates random string of length 10
func RandString10() string {
	return RandString(10)
}
