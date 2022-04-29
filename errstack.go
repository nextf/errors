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
