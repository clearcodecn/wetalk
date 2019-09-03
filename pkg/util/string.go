package util

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Md5(v string) string {
	m := md5.New()
	m.Write([]byte(v))
	data := m.Sum(nil)
	return fmt.Sprintf("%x", data)
}

func UUID() string {
	return uuid.NewV4().String()
}

const (
	ascii = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	num   = "0123456789"
)

func RandN(n int) string {
	x := make([]byte, n)
	for i := 0; i < n; i++ {
		x[i] = ascii[rand.Intn(len(ascii))]
	}
	return string(x)
}

func RandNumN(n int) string {
	x := make([]byte, n)
	for i := 0; i < n; i++ {
		x[i] = num[rand.Intn(len(num))]
	}
	return string(x)
}
