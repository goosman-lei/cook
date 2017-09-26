package io

import (
	"io"
	"io/ioutil"
)

func ReadAll_string(r io.Reader) (string, error) {
	s, e := ioutil.ReadAll(r)
	return string(s), e
}

func ReadN(r io.Reader, nExpected int) ([]byte, error) {
	var (
		buffer       []byte = make([]byte, nExpected)
		err          error
		nRead, tRead int = 0, 0
	)

	for nRead < nExpected {
		tRead, err = r.Read(buffer[nRead:])
		if err != nil {
			return nil, err
		}
		nRead += tRead
	}
	return buffer, nil

}

func ReadN_ToBuffer(r io.Reader, buffer []byte) error {
	var (
		err          error
		nRead, tRead int = 0, 0
		nExpected    int = len(buffer)
	)

	for nRead < nExpected {
		tRead, err = r.Read(buffer[nRead:])
		if err != nil {
			return err
		}
		nRead += tRead
	}
	return nil
}

func WriteN(w io.Writer, buffer []byte) error {
	var (
		err            error
		nWrite, tWrite int = 0, 0
		nExpected      int = len(buffer)
	)

	for nWrite < nExpected {
		tWrite, err = w.Write(buffer[nWrite:])
		if err != nil {
			return err
		}
		nWrite += tWrite
	}
	return nil

}
