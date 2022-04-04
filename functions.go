// This package is used to replace `errors` package of the standard library.
package errors

import (
	stderr "errors"
	"fmt"

	"github.com/nextf/errors/internal/pkgerr"
)

var (
	// Convergent the functions of standard library
	// func Is(err, target error) bool
	Is = stderr.Is

	// Convergent the functions of standard library
	// func As(err error, target interface{}) bool
	As = stderr.As

	// Convergent the functions of standard library
	// func Unwrap(err error) error
	Unwrap = stderr.Unwrap

	// Convergent the functions of pkg/errors library
	// WithMessage is a function that annotates err with a new message.
	// If err is nil, WithMessage returns nil.
	// func WithMessage(err error, message string) error
	WithMessage = pkgerr.WithMessage

	// Convergent the functions of pkg/errors library
	// WithMessagef is a function that annotates err with the format specifier.
	// If err is nil, WithMessagef returns nil.
	// func WithMessagef(err error, format string, args ...interface{}) error
	WithMessagef = pkgerr.WithMessagef
)

type StackTrace = pkgerr.StackTrace

func Match(err error, key interface{}) bool {
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
	for {
		if x, ok := err.(interface{ Code() string }); ok {
			return x.Code(), true
		}
		if err = Unwrap(err); err == nil {
			return "", false
		}
	}
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
func Wrap(err error, message string) error {
	return withStackIfAbsent(1, WithMessage(err, message))
}

// var Wrapf = pkgerr.Wrapf
func Wrapf(err error, format string, args ...interface{}) error {
	return withStackIfAbsent(1, WithMessage(err, fmt.Sprintf(format, args...)))
}

// var Errorf = pkgerr.Errorf
func Errorf(format string, args ...interface{}) error {
	return pkgerr.Errorf(1, format, args...)
}

// var NewWithStack = pkgerr.New
func NewWithStack(message string) error {
	return pkgerr.New(1, message)
}
