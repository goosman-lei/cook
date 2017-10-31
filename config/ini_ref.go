package config

import (
	"github.com/go-ini/ini"
	"regexp"
	"strings"
)

var (
	ref_fname                    = "ref.ini"
	ref_config map[string]string = make(map[string]string)

	pattern_ref = regexp.MustCompile("\\{\\$[-\\w]+\\.[-\\w]+\\}")
)

func Set_ref_fname(fname string) {
	ref_fname = fname
}

func Register_ref(key, value string) {
	ref_config[key] = value
}

func is_ref_file(fname string) bool {
	return fname == ref_fname
}

func init_ref_config(dname string) error {
	var (
		fname string = strings.TrimSuffix(dname, "/") + "/" + strings.TrimPrefix(ref_fname, "/")
		fp    *ini.File
		err   error
	)
	// refresh every time
	ref_config = make(map[string]string)

	// open as ini.File
	fp = ini.Empty()
	if err = fp.Append(fname); err != nil {
		return err
	}

	for _, sec := range fp.Sections() {
		for _, key := range sec.Keys() {
			ref_config[sec.Name()+"."+key.Name()] = key.Value()
		}
	}

	return nil
}

func replace_with_ref(content string) string {
	return pattern_ref.ReplaceAllStringFunc(content, func(matched string) string {
		if v, ok := ref_config[matched[2:len(matched)-1]]; ok {
			return v
		}
		return matched
	})
}
