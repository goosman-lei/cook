package zlib

import (
	"testing"
)

func Test_compress_and_decompress(t *testing.T) {
	s := "{\"id\":137668006515310592,\"tag_name\":\"\xe3\x82\x84\xe3\x81\x8d\xe3\x81\xa8\xe3\x82\x8a \xe7\x9a\x86\xe3\x81\xae\xe5\xae\xb6\",\"sense\":\"poi\",\"type\":\"point\",\"md5_tag_name\":\"066ec3fd998665098bee9ef0118db82a\",\"property\":\"{\\\"tag_id\\\":\\\"6352920\\\",\\\"tag_type\\\":\\\"custom_point\\\"}\",\"add_time\":1485686820,\"update_time\":1485686820}"

	z_s := Compress_string(s)
	unz_s := Decompress_string(z_s)

	if unz_s != s {
		t.Logf("compress or decompress error")
		t.Fail()
	}
}
