package util

import (
	"os"
	"path/filepath"
	"strings"
)

var execFile string
var execDir string
var workDir string

func init() {
	execFname, err := os.Executable()
	if err != nil {
		panic(err)
	}

	execFname, err = filepath.EvalSymlinks(execFname)
	if err != nil {
		panic(err)
	}

	workDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	execFile = execFname
	execDir = filepath.Dir(execFname)
}

func ExecFile() string {
	return execFile
}

func ExecDir() string {
	return execDir
}

func WorkDir() string {
	return workDir
}

func RuntimePath(path string) string {
	path = strings.Replace(path, "${ExecDir}", execDir, -1)
	path = strings.Replace(path, "${WorkDir}", workDir, -1)
	return path
}
