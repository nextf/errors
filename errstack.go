package errors

import (
	"fmt"
	"io"
	"strings"

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
			// Skip trace information
			if c.cause != nil {
				fmt.Fprintf(s, "%-v", c.cause)
			}
			break
		}
		// pretty print for log
		stacklog := strings.ReplaceAll(fmt.Sprintf("@callstack\n%v", c.stack), "\n", "\n\x20\x20\x20\x20")
		io.WriteString(s, stacklog)
		if s.Flag('+') && c.cause != nil {
			fmt.Fprintf(s, "\nCaused by: %+v", c.cause)
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
