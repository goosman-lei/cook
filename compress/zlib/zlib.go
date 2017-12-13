package zlib

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func Compress_string(s string) string {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write([]byte(s))
	w.Close()
	return buf.String()
}

func Decompress_string(s string) string {
	buf := bytes.NewBufferString(s)
	if r, err := zlib.NewReader(buf); err != nil {
		return ""
	} else if byts, err := ioutil.ReadAll(r); err != nil {
		return ""
	} else {
		return string(byts)
	}
}
