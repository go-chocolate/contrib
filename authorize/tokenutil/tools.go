package tokenutil

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

const (
	seed = "abcdef0123456789"
)

func randString(n int) string {
	var r = make([]byte, n)
	for i := 0; i < n; i++ {
		r[i] = seed[rand.Intn(len(seed))]
	}
	return string(r)
}

func toMd5(data []byte) string {
	b := md5.Sum(data)
	return hex.EncodeToString(b[:])
}
