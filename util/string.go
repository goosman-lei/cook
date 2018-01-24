package util

import (
	"bytes"
	"math/rand"
	"strings"
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

func Hump_to_underline(s string) string {
	b := new(bytes.Buffer)
	for i, c := range s {
		if c >= 64 && c <= 90 {
			if i != 0 {
				b.WriteRune('_')
			}
			b.WriteRune(c + 32)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func Equal_hump_to_underline(s1, s2 string) bool {
	return strings.EqualFold(Hump_to_underline(s1), Hump_to_underline(s2))
}
