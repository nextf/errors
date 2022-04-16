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

func New(code, message string) error {
	return &withErrCode{code, message, nil}
}

func Newf(code, format string, args ...interface{}) error {
	return &withErrCode{code, fmt.Sprintf(format, args...), nil}
}

func NewWithStack(code, message string) error {
	return &withErrCode{code, message, newErrorStack(1)}
}

func NewWithStackf(code, format string, args ...interface{}) error {
	return &withErrCode{code, fmt.Sprintf(format, args...), newErrorStack(1)}
}

func WithCode(cause error, code, message string) error {
	if cause == nil {
		return nil
	}
	if _, ok := GetCode(cause); ok {
		return cause
	}
	return &withErrCode{code, message, cause}
}

func WithCodePile(cause error, code, message string) error {
	if cause == nil {
		return nil
	}
	return &withErrCode{code, message, cause}
}

func HasErrorStack(err error) bool {
	if err == nil {
		return false
	}
	for {
		if _, ok := err.(interface{ CallersFrames() []stack.Frame }); ok {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func withStackIfAbsent(skip int, err error) error {
	if HasErrorStack(err) {
		return err
	} else {
		return withErrorStack(err, skip+1)
	}
}

// var WithStack = pkgerr.WithStack
func WithStack(cause error) error {
	if cause == nil {
		return nil
	}
	return withStackIfAbsent(1, cause)
}

func WithStackPile(err error) error {
	if err == nil {
		return nil
	}
	return withErrorStack(err, 1)
}

func Wrap(err error, code, message string) error {
	if err == nil {
		return nil
	}
	err = withStackIfAbsent(1, err)
	if _, ok := GetCode(err); !ok {
		err = &withErrCode{code, message, err}
	}
	return err
}

func Wrapf(err error, code, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = withStackIfAbsent(1, err)
	if _, ok := GetCode(err); !ok {
		err = &withErrCode{code, fmt.Sprintf(format, args...), err}
	}
	return err
}

func WrapPile(err error, code, message string) error {
	if err == nil {
		return nil
	}
	return &withErrCode{code, message, withErrorStack(err, 1)}
}

func WrapPilef(err error, code, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withErrCode{code, fmt.Sprintf(format, args...), withErrorStack(err, 1)}
}

// Deprecated: Too simple. Use errors.New instead.
func Message(message string) error {
	return stderr.New(message)
}

// Deprecated: Too simple. Use errors.Newf instead.
func Messagef(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Deprecated: Too simple. Use errors.Wrap instead.
func WithMessage(err error, message string) error {
	return &withMessage{message, withStackIfAbsent(1, err)}
}

// Deprecated: Too simple. Use errors.Wrapf instead.
func WithMessagef(err error, format string, args ...interface{}) error {
	return &withMessage{fmt.Sprintf(format, args...), withStackIfAbsent(1, err)}
}
