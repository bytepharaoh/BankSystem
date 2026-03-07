package util

import (
	"math/rand"
	"strings"
	"time"
)
const alphabet = "abcdefghijklmnopqrstuvwxyz"
func init(){
	rand.Seed(time.Now().UnixMicro())
}
// RandomInt generates a random numbers between min and max
func RandomInt(min , max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
// RandomString generates a random string of length n
func RandomString (n int) string {
	var sb strings.Builder
	k:= len(alphabet)
	for i := 0; i < n; i++ {
		c:=alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return  RandomString(6)
}
// RandomBalance generates a random amount of money
func RandomBalance() int64 {
	return RandomInt(1, 10000)
}

//RandomCurrency generate a random currency code
func RandomCurrency () string {
	currencies:= []string {USD , EGP ,RUB , EUR }
	n:=len(currencies)
	return currencies[rand.Intn(n)]
} 