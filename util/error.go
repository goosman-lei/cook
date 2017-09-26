package util

import (
	"io"
	"strings"
)

func Err_IsEof(err error) bool {
	return err == io.EOF
}

func Err_IsClosed(err error) bool {
	return strings.HasSuffix(err.Error(), "use of closed network connection")
}

func Err_IsTimeout(err error) bool {
	return strings.HasSuffix(err.Error(), "i/o timeout")
}

func Err_IsBroken(err error) bool {
	return strings.HasSuffix(err.Error(), "broken pipe")
}

func Err_NoSuchFileOrDir(err error) bool {
	return strings.HasSuffix(err.Error(), "no such file or directory")
}
