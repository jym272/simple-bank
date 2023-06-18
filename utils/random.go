package utils

import (
	"math/rand"
	"time"
)

type Random struct {
	rnm *rand.Rand
}

func getRand() *Random {
	r := &Random{}
	r.init()
	return r
}

func (r *Random) init() {
	r.rnm = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (r *Random) RandomInt(min, max int64) int64 {
	return min + r.rnm.Int63n(max-min+1)
}

func (r *Random) RandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.rnm.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomOwner() string {
	return getRand().RandomString(6)
}

func RandomMoney() int64 {
	return getRand().RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[getRand().RandomInt(0, int64(n-1))]
}
