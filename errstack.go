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

	"github.com/nextf/errors/stack"
)

const (
	maxDepth int = 32
)

type errorStack struct {
	stack stack.CallStack
	cause error
}

func (c *errorStack) StackTrace() []stack.Frame {
	return c.stack.StackTrace()
}

func (c *errorStack) Error() string {
	return c.cause.Error()
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (c *errorStack) Unwrap() error {
	return c.cause
}

func (c *errorStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('-') {
			// Skip stack trace
			if c.cause != nil {
				fmt.Fprintf(s, "%-v", c.cause)
			}
			break
		}
		// Print stack trace
		width, hasWidth := s.Width()
		formatCallStack := "@callstack\n%+v"
		if hasWidth {
			formatCallStack = fmt.Sprintf("@callstack\n%%+%dv", width)
		}
		fmt.Fprintf(s, formatCallStack, c.stack)
		if s.Flag('+') && c.cause != nil {
			formatCause := "\nCaused by: %+v"
			if hasWidth {
				formatCause = fmt.Sprintf("\nCaused by: %%+%dv", width)
			}
			fmt.Fprintf(s, formatCause, c.cause)
		}
	case 's':
		if c.cause != nil {
			fmt.Fprintf(s, "%s", c.cause)
		}
	case 'q':
		if c.cause != nil {
			fmt.Fprintf(s, "%q", c.cause)
		}
	}
}

func newErrorStack(skip int) error {
	return &errorStack{stack.RecordCallStack(skip+1, maxDepth), nil}
}

func withErrorStack(err error, skip int) error {
	return &errorStack{stack.RecordCallStack(skip+1, maxDepth), err}
}
