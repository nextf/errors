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
	// err2 := errors.WithStack(err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	// fmt.Printf("%+v\n", err2)
	// Output:
	// whoops
	// [ERR001] whoops
}

func ExampleNew_printf2() {
	err := errors.NewWithStack("ERR001", "whoops")
	fmt.Printf("%+v\n", err)
	// Output:
	// [ERR001] whoops
	// Caused by: @callstack
	// github.com/nextf/errors_test.ExampleNew_printf2(example_test.go:29)
	// testing.runExample(run_example.go:64)
	// testing.runExamples(example.go:44)
	// testing.(*M).Run(testing.go:1505)
	// main.main(_testmain.go:75)
	// runtime.main(proc.go:255)
	// runtime.goexit(asm_amd64.s:1581)
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

func openNotExistsFile2() (*os.File, error) {
	f, err := os.Open("/not_exists_file.txt")
	if err != nil {
		return nil, errors.Wrap(err, "ERR404", "File not found")
	}
	return f, nil
}

func ExampleWrap() {
	_, err := openNotExistsFile2()
	fmt.Printf("%+v", err)
	// Output:
	// [ERR404] File not found
	// Caused by: @callstack
	// github.com/nextf/errors_test.openNotExistsFile2(example_test.go:61)
	// github.com/nextf/errors_test.ExampleWrap(example_test.go:67)
	// testing.runExample(run_example.go:64)
	// testing.runExamples(example.go:44)
	// testing.(*M).Run(testing.go:1505)
	// main.main(_testmain.go:75)
	// runtime.main(proc.go:255)
	// runtime.goexit(asm_amd64.s:1581)
	// Caused by: open /not_exists_file.txt: The system cannot find the file specified.
}
