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
	"fmt"
	"io"
)

type withErrCode struct {
	code    string
	message string
	cause   error
}

func (c *withErrCode) Error() string {
	return c.message
}

func (c *withErrCode) Code() string {
	return c.code
}

func (c *withErrCode) Match(key interface{}) bool {
	if key == nil {
		return false
	}
	switch x := key.(type) {
	case string:
		return c.code == x
	case interface{ MatchString(s string) bool }:
		return x.MatchString(c.code)
	}
	return false
}

func (c *withErrCode) Unwrap() error {
	return c.cause
}

func (c *withErrCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "[%s] %s", c.code, c.message)
		if s.Flag('+') && c.cause != nil {
			formatCause := "\nCaused by: %+v"
			if width, ok := s.Width(); ok {
				formatCause = fmt.Sprintf("\nCaused by: %%+%dv", width)
			}
			fmt.Fprintf(s, formatCause, c.cause)
		}
	case 's':
		io.WriteString(s, c.message)
	case 'q':
		fmt.Fprintf(s, "%q", c.message)
	}
}
