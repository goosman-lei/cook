package util

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	RAND_STR_MAPPING     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	RAND_STR_MAPPING_LEN = 62
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Rand_str(len int) string {
	buf := new(bytes.Buffer)

	for ; len > 0; len-- {
		buf.WriteByte(RAND_STR_MAPPING[rand.Intn(RAND_STR_MAPPING_LEN)])
	}

	return buf.String()
}
