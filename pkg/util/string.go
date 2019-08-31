package util

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func Md5(v string) string {
	m := md5.New()
	m.Write([]byte(v))
	data := m.Sum(nil)
	return fmt.Sprintf("%x", data)
}

func UUID() string {
	return uuid.NewV4().String()
}
