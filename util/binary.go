package util

import (
	"encoding/binary"
	"unsafe"
)

var (
	MachineEndian binary.ByteOrder
)

const (
	INT_SIZE int = int(unsafe.Sizeof(0))
)

func init() {
	var i int = 0x1
	if d := (*[INT_SIZE]byte)(unsafe.Pointer(&i)); d[0] == 0 {
		MachineEndian = binary.BigEndian
	} else {
		MachineEndian = binary.LittleEndian
	}
}
