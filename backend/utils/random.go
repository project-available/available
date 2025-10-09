package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numberic = "0123456789"


// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomStringNumber(n int) string {
	var sb strings.Builder
	k := len(numberic)

	for i := 0; i < n; i++ {
		c := numberic[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomRole() string {
	roles := []string{"user", "admin"}
	return roles[rand.Intn(len(roles))]
}

func RandomEmail() string {
	email := RandomString(15) + "@hcmut.edu.vn"
	return email
}

func RandomPhone() string {
	return RandomStringNumber(10)
}

func RandomStudentID() string {
	return RandomStringNumber(7)
}
