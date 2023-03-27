package utils

import "math/rand"

func RandomString(length int) string {
	s := make([]byte, length)

	for i := 0; i < length; i++ {
		s[i] = 'a' + byte(rand.Intn(26))
	}

	return string(s)
}
