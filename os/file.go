package os

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"os"
	"path/filepath"
)

func OpenFileWithMkdir(path string, flag int, perm os.FileMode) (*os.File, error) {
	var (
		fp  *os.File
		err error
	)
	if fp, err = os.OpenFile(path, flag, perm); err == nil {
		return fp, err
	}

	if !cook_util.Err_NoSuchFileOrDir(err) {
		return nil, err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(path, flag, perm)
}

func IsRegularFileExist(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		return false
	}

	fm := fd.Mode()

	if fm.IsRegular() {
		return true
	}
	return false
}
