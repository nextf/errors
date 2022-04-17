package errors

import (
	"fmt"
	"io"
)

type errorMessage struct {
	message string
	cause   error
}

func (c *errorMessage) Error() string {
	return c.message
}

func (c *errorMessage) Unwrap() error {
	return c.cause
}

func (c *errorMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, c.message)
		if s.Flag('+') && c.cause != nil {
			fmt.Fprintf(s, "\nCaused by: %+v", c.cause)
		}
	case 's':
		io.WriteString(s, c.message)
	case 'q':
		fmt.Fprintf(s, "%q", c.message)
	}
}
