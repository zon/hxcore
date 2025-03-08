package test

import (
	"github.com/zon/hxcore"
)

func RandomEmail() string {
	return hxcore.RandomString(8) + "@example.com"
}

func RandomPassword() string {
	return hxcore.RandomString(16)
}
