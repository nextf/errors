package errors

import (
	"fmt"
	"io"

	pkgerr "github.com/pkg/errors"
)

type codedError struct {
	code    string
	message string
	error
}

func (c *codedError) Error() string {
	return c.message
}

func (c *codedError) Code() string {
	return c.code
}

// func (ce *codedError) Is(err error) bool {
// 	var coder interface{ Code() string }
// 	if As(err, &coder) {
// 		return ce.Code() == coder.Code()
// 	}
// 	return false
// }

func (c *codedError) Format(s fmt.State, verb rune) {
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
	return &codedError{code, message, nil}
}

func CodErrWithStack(code string, message string) error {
	ce := &codedError{code, message, nil}
	return pkgerr.WithStack(ce)
}

func CodErrWrap(code string, message string, cause error) error {
	if cause == nil {
		return nil
	}
	return &codedError{code, message, WithStack(cause)}
}
