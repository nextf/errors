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

type errorMessage struct {
	message string
	cause   error
}

func (c *errorMessage) Error() string {
	return c.message
}

// Unwrap provides compatibility for Go 1.13 error chains.
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
