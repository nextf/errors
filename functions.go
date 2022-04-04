// This package is used to replace `errors` package of the standard library.
package errors

import (
	stderr "errors"
	"fmt"

	"github.com/nextf/errors/internal/pkgerr"
)

var (
	Is     = stderr.Is
	As     = stderr.As
	Unwrap = stderr.Unwrap

	Errorf       = pkgerr.Errorf
	NewSimple    = pkgerr.New
	WithMessage  = pkgerr.WithMessage
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

// var WithStack = pkgerr.WithStack
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if HasStackTrace(err) {
		return err
	} else {
		return pkgerr.WithStack(pkgerr.SkipPkgErr+1, err)
	}
}

// var Wrap = pkgerr.Wrap
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return WithStack(WithMessage(err, message))
}

// var Wrapf = pkgerr.Wrapf
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return Wrap(err, fmt.Sprintf(format, args...))
}
