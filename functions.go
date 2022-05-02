// Copyright 2022 Zhiwen<zhiwen.t@outlook.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	stderr "errors"
	"fmt"

	"github.com/nextf/errors/stack"
)

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library.
func Is(err, target error) bool {
	return stderr.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
func As(err error, target interface{}) bool {
	return stderr.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderr.Unwrap(err)
}

// Match reports whether any error in err's chain matches key.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error if it implements a method Match(key) bool such that Match(target)
// returns true.
//
// An error type might provide an Match method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Match(key interface{}) bool { return m.code == key }
//
// then Match(MyError{code:"ERR001"}, "ERR001") returns true.
func Match(err error, target interface{}) bool {
	if err == nil {
		return false
	}
	for {
		if x, ok := err.(interface{ Match(interface{}) bool }); ok && x.Match(target) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// GetCode finds the first error in err's chain that implements a method Code() string,
// and if so, returns the code extracted from error and the boolean is true.
// Otherwise the returned value will be empty and the boolean will be false.
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

func HasStackTrace(err error) bool {
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
	if HasStackTrace(err) {
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

// New returns an error with the supplied message.
// If a message begins with code enclosed in [], that code is considered an error code.
func New(message string) error {
	return ConstError(message)
}

// Errorf formats according to a format specifier and returns the string as a value that satisfies error.
// If a string begins with code enclosed in [], that code is considered an error code.
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
