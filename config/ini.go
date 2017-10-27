package config

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"strings"
)

const (
	INI_FILE_SUFFIX = ".ini"
)

var (
	ErrSectionNotExists = errors.New("section not exists")
)

func Ini_open_dir(dname string) (*ini.File, error) {
	dp, err := os.Open(dname)
	if err != nil {
		return nil, err
	}

	fnames, err := dp.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	init_ref_config(dname)

	fp := ini.Empty()

	for _, fname := range fnames {
		if !strings.HasSuffix(fname, INI_FILE_SUFFIX) {
			continue
		}
		if is_ref_file(fname) {
			continue
		}

		if content, err := ioutil.ReadFile(dname + "/" + fname); err != nil {
			return nil, err
		} else if err := fp.Append([]byte(replace_with_ref(string(content)))); err != nil {
			return nil, err
		}
	}

	return fp, nil
}

func Ini_direct_get_key(sec *ini.Section, node, key string) *ini.Key {
	if k, err := sec.GetKey(key); err == nil {
		return k
	}
	panic(fmt.Errorf("Config error: [%s#%s] not exists", node, key))
	return nil
}
func Ini_inherit_get_key(sec, psec *ini.Section, node, key string) *ini.Key {
	if k, err := sec.GetKey(key); err == nil {
		return k
	}
	if k, err := psec.GetKey(key); err == nil {
		return k
	}
	panic(fmt.Errorf("Config error: [%s#%s] not exists", node, key))
	return nil
}
