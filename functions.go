// This package is used to replace `errors` package of the standard library.
package errors

import (
	stderr "errors"
	"fmt"

	"github.com/nextf/errors/internal/pkgerr"
)

type StackTrace = pkgerr.StackTrace

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
	return pkgerr.WithStack(1, &withErrCode{code, message, nil})
}

func NewWithStackf(code, format string, args ...interface{}) error {
	return pkgerr.WithStack(1, &withErrCode{code, fmt.Sprintf(format, args...), nil})
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

func WithCodeMass(cause error, code, message string) error {
	if cause == nil {
		return nil
	}
	return &withErrCode{code, message, cause}
}

func HasStackTrace(err error) bool {
	if err == nil {
		return false
	}
	for {
		if _, ok := err.(interface{ StackTrace() StackTrace }); ok {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func withStackIfAbsent(skipCounter int, err error) error {
	if HasStackTrace(err) {
		return err
	} else {
		return pkgerr.WithStack(skipCounter+1, err)
	}
}

// var WithStack = pkgerr.WithStack
func WithStack(cause error) error {
	if cause == nil {
		return nil
	}
	return withStackIfAbsent(1, cause)
}

func WithStackMass(cause error) error {
	if cause == nil {
		return nil
	}
	return pkgerr.WithStack(1, cause)
}

func Wrap(cause error, code, message string) error {
	if cause == nil {
		return nil
	}
	if _, ok := GetCode(cause); ok {
		return withStackIfAbsent(1, cause)
	}
	return withStackIfAbsent(1, &withErrCode{code, message, cause})
}

func Wrapf(cause error, code, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}
	if _, ok := GetCode(cause); ok {
		return withStackIfAbsent(1, cause)
	}
	return withStackIfAbsent(1, &withErrCode{code, fmt.Sprintf(format, args...), cause})
}

func WrapMass(cause error, code, message string) error {
	if cause == nil {
		return nil
	}
	return pkgerr.WithStack(1, &withErrCode{code, message, cause})
}

func WrapMassf(cause error, code, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}
	return pkgerr.WithStack(1, &withErrCode{code, fmt.Sprintf(format, args...), cause})
}

// Deprecated: Too simple. Use errors.New instead.
func Message(message string) error {
	return stderr.New(message)
}

// Deprecated: Too simple. Use errors.New instead.
func Messagef(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Deprecated: Too simple. Use errors.Wrap instead.
func WithMessage(err error, message string) error {
	return withStackIfAbsent(1, pkgerr.WithMessage(err, message))
}

// Deprecated: Too simple. Use errors.Wrapf instead.
func WithMessagef(err error, format string, args ...interface{}) error {
	return withStackIfAbsent(1, pkgerr.WithMessagef(err, format, args...))
}
