package ginx

import "strconv"

type Code int

func (c Code) String() string {
	return strconv.Itoa(int(c))
}

type E struct {
	message string
	code    Code
	exist   bool
}

func (e E) Error() string {
	return e.message
}

func (e E) Code() string {
	return e.code.String()
}

func (e E) CodeInt() int {
	return int(e.code)
}

func (e E) Empty() bool {
	return e.exist
}

func EmptyError() E {
	return E{exist: false}
}

func NewError(msg string) E {
	return E{code: Error, message: msg, exist: true}
}

func NewWithCode(msg string, code Code) E {
	return E{code: code, message: msg, exist: true}
}
