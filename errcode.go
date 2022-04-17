package errors

import (
	"fmt"
	"io"
)

type withErrCode struct {
	code    string
	message string
	cause   error
}

func (c *withErrCode) Error() string {
	return c.message
}

func (c *withErrCode) Code() string {
	return c.code
}

func (c *withErrCode) Match(key interface{}) bool {
	switch x := key.(type) {
	case string:
		return c.code == x
	case interface{ MatchString(s string) bool }:
		return x.MatchString(c.code)
	}
	return false
}

func (c *withErrCode) Unwrap() error {
	return c.cause
}

// func (ce *codedError) Is(err error) bool {
// 	var coder interface{ Code() string }
// 	if As(err, &coder) {
// 		return ce.Code() == coder.Code()
// 	}
// 	return false
// }

func (c *withErrCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "[%s] %s", c.code, c.message)
		if s.Flag('+') && c.cause != nil {
			fmt.Fprintf(s, "\nCaused by: %+v", c.cause)
		}
	case 's':
		io.WriteString(s, c.message)
	case 'q':
		fmt.Fprintf(s, "%q", c.message)
	}
}
