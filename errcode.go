package errors

import (
	"fmt"
	"io"

	pkgerr "github.com/nextf/errors/internal/pkgerr"
)

type withErrCode struct {
	code    string
	message string
	error
}

func (c *withErrCode) Error() string {
	return c.message
}

func (c *withErrCode) Code() string {
	return c.code
}

func (c *withErrCode) Match(key interface{}) bool {
	if x, ok := key.(string); ok {
		return c.code == x
	}
	return false
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
		fmt.Fprintf(s, "[%s] %s", c.Code(), c.Error())
	case 's':
		io.WriteString(s, c.message)
	case 'q':
		fmt.Fprintf(s, "%q", c.message)
	}
}

func ErrCode(code string, message string) error {
	return &withErrCode{code, message, nil}
}

func ErrCodeWithStack(code string, message string) error {
	ce := &withErrCode{code, message, nil}
	return pkgerr.WithStack(pkgerr.SkipPkgErr+1, ce)
}

func ErrCodeWrap(code string, message string, cause error) error {
	if cause == nil {
		return nil
	}
	return &withErrCode{code, message, WithStack(cause)}
}
