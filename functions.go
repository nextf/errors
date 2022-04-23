// This package is used to replace `errors` package of the standard library.
package errors

import (
	stderr "errors"
	"fmt"

	"github.com/nextf/errors/stack"
)

func Is(err, target error) bool {
	return stderr.Is(err, target)
}

func As(err error, target interface{}) bool {
	return stderr.As(err, target)
}

func Unwrap(err error) error {
	return stderr.Unwrap(err)
}

func Match(err error, key interface{}) bool {
	if err == nil {
		return false
	}
	for {
		if x, ok := err.(interface{ Match(interface{}) bool }); ok && x.Match(key) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func GetCode(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	for {
		if x, ok := err.(interface{ Code() string }); ok {
			return x.Code(), true
		}
		if err = Unwrap(err); err == nil {
			return "", false
		}
	}
}

func ErrCode(code, message string) error {
	return &withErrCode{code, message, nil}
}

func ErrCodef(code, format string, args ...interface{}) error {
	return &withErrCode{code, fmt.Sprintf(format, args...), nil}
}

func TraceableErrCode(code, message string) error {
	return &withErrCode{code, message, newErrorStack(1)}
}

func TraceableErrCodef(code, format string, args ...interface{}) error {
	return &withErrCode{code, fmt.Sprintf(format, args...), newErrorStack(1)}
}

// func WithNodupErrCode(err error, code, message string) error {
// 	if err == nil {
// 		return nil
// 	}
// 	if x, ok := GetCode(err); ok && code == x {
// 		return err
// 	}
// 	return &withErrCode{code, message, err}
// }

func WithErrCode(err error, code, message string) error {
	if err == nil {
		return nil
	}
	return &withErrCode{code, message, err}
}

func WithErrCodef(err error, code, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withErrCode{code, fmt.Sprintf(format, args...), err}
}

func HasErrorStack(err error) bool {
	if err == nil {
		return false
	}
	for {
		if _, ok := err.(interface{ StackTrace() []stack.Frame }); ok {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func withStackIfAbsent(err error, skip int) error {
	if HasErrorStack(err) {
		return err
	} else {
		return withErrorStack(err, skip+1)
	}
}

func TraceNodup(err error) error {
	if err == nil {
		return nil
	}
	return withStackIfAbsent(err, 1)
}

func Trace(err error) error {
	if err == nil {
		return nil
	}
	return withErrorStack(err, 1)
}

func WrapNodup(err error, code, message string) error {
	if err == nil {
		return nil
	}
	err = withStackIfAbsent(err, 1)
	return &withErrCode{code, message, err}
}

func WrapNodupf(err error, code, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = withStackIfAbsent(err, 1)
	return &withErrCode{code, fmt.Sprintf(format, args...), err}
}

func Wrap(err error, code, message string) error {
	if err == nil {
		return nil
	}
	return &withErrCode{code, message, withErrorStack(err, 1)}
}

func Wrapf(err error, code, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withErrCode{code, fmt.Sprintf(format, args...), withErrorStack(err, 1)}
}

func New(message string) error {
	return ConstError(message)
}

func Errorf(format string, args ...interface{}) error {
	return ConstError(fmt.Sprintf(format, args...))
}

// Deprecated: Too simple. Use errors.Wrap instead.
func TraceMessage(err error, message string) error {
	return &errorMessage{message, withStackIfAbsent(err, 1)}
}

// Deprecated: Too simple. Use errors.Wrapf instead.
func TraceMessagef(err error, format string, args ...interface{}) error {
	return &errorMessage{fmt.Sprintf(format, args...), withStackIfAbsent(err, 1)}
}
