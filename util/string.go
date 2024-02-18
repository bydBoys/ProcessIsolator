package util

import (
	"math/rand"
	"time"
)

var randPool *rand.Rand

func init() {
	randPool = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandStringBytes(n int) string {
	letterBytes := "1234567890"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[randPool.Intn(len(letterBytes))]
	}
	return string(b)
}
