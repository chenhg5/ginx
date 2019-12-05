package ginx

type Config interface {
	Debug() bool
	Production() bool
}
