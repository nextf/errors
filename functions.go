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

func ErrCode(code, message string) error {
	return &withErrCode{code, message, nil}
}

func ErrCodeWithStack(code, message string) error {
	return pkgerr.WithStack(1, &withErrCode{code, message, nil})
}

func WithErrCode(code, message string, cause error) error {
	if cause == nil {
		return nil
	}
	return withStackIfAbsent(1, &withErrCode{code, message, cause})
}

func WithErrCodeIfAbsent(code, message string, cause error) error {
	if cause == nil {
		return nil
	}
	if _, ok := GetCode(cause); ok {
		return withStackIfAbsent(1, cause)
	}
	return withStackIfAbsent(1, &withErrCode{code, message, cause})
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
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if HasStackTrace(err) {
		return err
	} else {
		return withStackIfAbsent(1, err)
	}
}

// var Wrap = pkgerr.Wrap
func WithMessage(err error, message string) error {
	return withStackIfAbsent(1, pkgerr.WithMessage(err, message))
}

// var Wrapf = pkgerr.Wrapf
func WithMessagef(err error, format string, args ...interface{}) error {
	return withStackIfAbsent(1, pkgerr.WithMessagef(err, format, args...))
}

// var Errorf = fmt.Errorf
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// var Stack = pkgerr.New
func Stack(message string) error {
	return pkgerr.New(1, message)
}

func Stackf(format string, args ...interface{}) error {
	return pkgerr.Errorf(1, format, args...)
}
