package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecFileAndExecDir(t *testing.T) {
	var (
		absPath string
		err     error
	)
	absPath, err = filepath.Abs(os.Args[0])
	if err != nil {
		t.Logf("get test file runtime absolute path failed: %s", err)
		t.Fail()
	}
	absPath, err = filepath.EvalSymlinks(absPath)
	if err != nil {
		t.Logf("eval test file runtime abs path symlinks failed: %s", err)
		t.Fail()
	}

	if ExecFile() != absPath {
		t.Logf("wrong file: ExecFile = %s, TestingFile = %s", ExecFile(), absPath)
		t.Fail()
	}

	if ExecDir() != filepath.Dir(absPath) {
		t.Logf("wrong file: ExecDir = %s, TestingDir = %s", ExecDir(), filepath.Dir(absPath))
		t.Fail()
	}

}
