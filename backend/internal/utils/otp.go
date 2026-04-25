package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {
	if length <= 0 {
		length = 6
	}
	rand.Seed(time.Now().UnixNano())
	max := 1
	for i := 0; i < length; i++ {
		max *= 10
	}
	n := rand.Intn(max)
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n)
}
