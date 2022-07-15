package utils

import (
	"math/rand"
	"time"
)

const rayBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
const (
	rayIdxBits = 6                 // 6 bits to represent a ray index
	rayIdxMask = 1<<rayIdxBits - 1 // All 1-bits, as many as rayIdxBits
	rayIdxMax  = 63 / rayIdxBits   // # of ray indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func GetRayByte(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), rayIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), rayIdxMax
		}
		if idx := int(cache & rayIdxMask); idx < len(rayBytes) {
			b[i] = rayBytes[idx]
			i--
		}
		cache >>= rayIdxBits
		remain--
	}

	return b
}

func GetRay() string {
	return string(GetRayByte(16))
}
