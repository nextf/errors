package stack_test

import (
	"fmt"

	"github.com/nextf/errors/stack"
)

func loopCall(time uint) stack.CallStack {
	if time > 1 {
		time--
		return loopCall(time)
	}
	return stack.RecordCallStack(0, 5)
}

func ExampleNewErrorStack() {
	stack := loopCall(2)
	fmt.Printf("%v", stack)
	// Output:
	// github.com/nextf/errors/stack_test.loopCall(stack_test.go:14)
	// github.com/nextf/errors/stack_test.loopCall(stack_test.go:12)
	// github.com/nextf/errors/stack_test.ExampleNewErrorStack(stack_test.go:18)
	// testing.runExample(run_example.go:64)
	// testing.runExamples(example.go:44)
}

func ExampleNewErrorStack_limited() {
	stack := loopCall(2)
	fmt.Printf("%3v", stack)
	// Output:
	// github.com/nextf/errors/stack_test.loopCall(stack_test.go:14)
	// github.com/nextf/errors/stack_test.loopCall(stack_test.go:12)
	// github.com/nextf/errors/stack_test.ExampleNewErrorStack_limited(stack_test.go:29)
	// ...(more:2)
}

func ExampleNewErrorStack_indent() {
	stack := loopCall(2)
	fmt.Printf("%+3v", stack)
	// Output:
	//     github.com/nextf/errors/stack_test.loopCall(stack_test.go:14)
	//     github.com/nextf/errors/stack_test.loopCall(stack_test.go:12)
	//     github.com/nextf/errors/stack_test.ExampleNewErrorStack_indent(stack_test.go:39)
	//     ...(more:2)
}
