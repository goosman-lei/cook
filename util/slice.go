package util

import (
	"strconv"
)

func Slice_int_to_string(s []int) []string {
	sr := make([]string, len(s))
	for i, v := range s {
		sr[i] = strconv.FormatInt(int64(v), 10)
	}
	return sr
}
