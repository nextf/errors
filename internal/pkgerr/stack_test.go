package pkgerr

import (
	"fmt"
	"runtime"
	"testing"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		initpc,
		"%s",
		"stack_test.go",
	}, {
		initpc,
		"%+s",
		"github.com/nextf/errors/internal/pkgerr.init\n" +
			"\t.+/nextf/errors/internal/pkgerr/stack_test.go",
	}, {
		0,
		"%s",
		"unknown",
	}, {
		0,
		"%+s",
		"unknown",
	}, {
		initpc,
		"%d",
		"9",
	}, {
		0,
		"%d",
		"0",
	}, {
		initpc,
		"%n",
		"init",
	}, {
		func() Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`\(\*X\).ptr`,
	}, {
		func() Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		0,
		"%n",
		"",
	}, {
		initpc,
		"%v",
		"stack_test.go:9",
	}, {
		initpc,
		"%+v",
		"github.com/nextf/errors/internal/pkgerr.init\n" +
			"\t.+/nextf/errors/internal/pkgerr/stack_test.go:9",
	}, {
		0,
		"%v",
		"unknown:0",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.Frame, tt.format, tt.want)
	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/nextf/errors/internal/pkgerr.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want []string
	}{{
		New(SkipPkgErr, "ooh"), []string{
			"github.com/nextf/errors/internal/pkgerr.TestStackTrace\n" +
				"\t.+/nextf/errors/internal/pkgerr/stack_test.go:121",
		},
	}, {
		Wrap(SkipPkgErr, New(SkipPkgErr, "ooh"), "ahh"), []string{
			"github.com/nextf/errors/internal/pkgerr.TestStackTrace\n" +
				"\t.+/nextf/errors/internal/pkgerr/stack_test.go:126", // this is the stack of Wrap, not New
		},
	}, {
		Cause(Wrap(SkipPkgErr, New(SkipPkgErr, "ooh"), "ahh")), []string{
			"github.com/nextf/errors/internal/pkgerr.TestStackTrace\n" +
				"\t.+/nextf/errors/internal/pkgerr/stack_test.go:131", // this is the stack of New
		},
	}, {
		func() error { return New(SkipPkgErr, "ooh") }(), []string{
			`github.com/nextf/errors/internal/pkgerr.TestStackTrace.func1` +
				"\n\t.+/nextf/errors/internal/pkgerr/stack_test.go:136", // this is the stack of New
			"github.com/nextf/errors/internal/pkgerr.TestStackTrace\n" +
				"\t.+/nextf/errors/internal/pkgerr/stack_test.go:136", // this is the stack of New's caller
		},
	}, {
		Cause(func() error {
			return func() error {
				return Errorf(SkipPkgErr, "hello %s", fmt.Sprintf("world: %s", "ooh"))
			}()
		}()), []string{
			`github.com/nextf/errors/internal/pkgerr.TestStackTrace.func2.1` +
				"\n\t.+/nextf/errors/internal/pkgerr/stack_test.go:145", // this is the stack of Errorf
			`github.com/nextf/errors/internal/pkgerr.TestStackTrace.func2` +
				"\n\t.+/nextf/errors/internal/pkgerr/stack_test.go:146", // this is the stack of Errorf's caller
			"github.com/nextf/errors/internal/pkgerr.TestStackTrace\n" +
				"\t.+/nextf/errors/internal/pkgerr/stack_test.go:147", // this is the stack of Errorf's caller's caller
		},
	}}
	for i, tt := range tests {
		x, ok := tt.err.(interface {
			StackTrace() StackTrace
		})
		if !ok {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := x.StackTrace()
		for j, want := range tt.want {
			testFormatRegexp(t, i, st[j], "%+v", want)
		}
	}
}

func stackTrace() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st stack = pcs[0:n]
	return st.StackTrace()
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]pkgerr.Frame\(nil\)`,
	}, {
		make(StackTrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%+v",
		"",
	}, {
		make(StackTrace, 0),
		"%#v",
		`\[\]pkgerr.Frame{}`,
	}, {
		stackTrace()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stackTrace()[:2],
		"%v",
		`\[stack_test.go:174 stack_test.go:221\]`,
	}, {
		stackTrace()[:2],
		"%+v",
		"\n" +
			"github.com/nextf/errors/internal/pkgerr.stackTrace\n" +
			"\t.+/nextf/errors/internal/pkgerr/stack_test.go:174\n" +
			"github.com/nextf/errors/internal/pkgerr.TestStackTraceFormat\n" +
			"\t.+/nextf/errors/internal/pkgerr/stack_test.go:225",
	}, {
		stackTrace()[:2],
		"%#v",
		`\[\]pkgerr.Frame{stack_test.go:174, stack_test.go:233}`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

// a version of runtime.Caller that returns a Frame, not a uintptr.
func caller() Frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return Frame(frame.PC)
}