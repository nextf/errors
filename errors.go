package errors

import (
	"fmt"
	"io"

	pkgerr "github.com/pkg/errors"
)

type withCodeAndMsg struct {
	code    string
	message string
	error
}

func (c *withCodeAndMsg) Error() string {
	return c.message
}

func (c *withCodeAndMsg) Code() string {
	return c.code
}

func (c *withCodeAndMsg) Match(key interface{}) bool {
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

func (c *withCodeAndMsg) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "[%s] %s", c.Code(), c.Error())
	case 's':
		io.WriteString(s, c.message)
	case 'q':
		fmt.Fprintf(s, "%q", c.message)
	}
}

func CodErr(code string, message string) error {
	return &withCodeAndMsg{code, message, nil}
}

func CodErrWithStack(code string, message string) error {
	ce := &withCodeAndMsg{code, message, nil}
	return pkgerr.WithStack(ce)
}

func CodErrWrap(code string, message string, cause error) error {
	if cause == nil {
		return nil
	}
	return &withCodeAndMsg{code, message, WithStack(cause)}
}
