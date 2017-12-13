package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Compress_string(s string) string {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if _, err := zw.Write([]byte(s)); err != nil {
		return ""
	} else if err := zw.Close(); err != nil {
		return ""
	} else {
		return buf.String()
	}
}

func Decompress_string(s string) string {
	buf := bytes.NewBufferString(s)
	if zr, err := gzip.NewReader(buf); err != nil {
		return ""
	} else if byts, err := ioutil.ReadAll(zr); err != nil {
		return ""
	} else {
		return string(byts)
	}
}
