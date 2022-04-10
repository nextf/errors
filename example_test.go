package errors_test

import (
	"fmt"
	"os"

	"github.com/nextf/errors"
)

func ExampleNew() {
	err := errors.New("ERR001", "whoops")
	fmt.Println(err)
	// Output:
	// [ERR001] whoops
}

func ExampleNew_printf1() {
	err := errors.New("ERR001", "whoops")
	err2 := errors.WithStack(err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err2)
	// Output:
	// whoops
	// [ERR001] whoops
}

func ExampleNew_printf2() {
	err := errors.New("ERR001", "whoops")
	err2 := errors.WithStack(err)
	fmt.Printf("%+v\n", err2)
	// Simple output:
	// [ERR001] whoops
	// github.com/nextf/errors_test.ExampleNew_printf
	// 	/src/errors/example_test.go:18
	// testing.runExample
	// 	/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:71
	// runtime.main
	// 	/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/go/src/runtime/asm_amd64.s:1581
}

func openNotExistsFile() (*os.File, error) {
	f, err := os.Open("/not_exists_file.txt")
	if err != nil {
		return nil, errors.WithCode(err, "ERR404", "File not found")
	}
	return f, nil
}
func ExampleWithCode() {
	_, err := openNotExistsFile()
	fmt.Printf("%+v", err)
	// Output:
	// [ERR404] File not found
	// Caused by: open /not_exists_file.txt: The system cannot find the file specified.
}
