package orm

var (
	InDebug bool = false
)

func DebugOn() {
	InDebug = true
}

func DebugOff() {
	InDebug = false
}
