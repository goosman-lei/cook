package util

import (
	"os"
)

var (
	Hostname string
)

func init() {
	Hostname, _ = os.Hostname()
}
