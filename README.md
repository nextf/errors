# errors
[![Go Report Card](https://goreportcard.com/badge/github.com/nextf/errors)](https://goreportcard.com/report/github.com/nextf/errors)
[![GoDoc](https://godoc.org/github.com/nextf/errors?status.svg)](http://godoc.org/github.com/nextf/errors)

This go package provides error handling primitives based on error codes.

# Introduction
Error coding is an ancient method of error orchestration. It's not elegant, but it works.
When errors have codes, they are easy to transport, tag, locate, and classify.
It would be more conducive to aspect oriented programming.

# Usage
Full documentation is available on [godoc](http://godoc.org/github.com/nextf/errors), but here's a simple example:
## Sentinel Error
Declare sentinel errors with constants to avoid tampering.
```go
const (
	ErrNotFoundOrders = errors.ConstError("[NF_BIS_Order] Not found orders")
	ErrDbAccessDeny   = errors.ConstError("[AD_TEC_DbConnect] Database access denied")
)
```
## Adding context to an error
The errors.Wrap function returns a new error that adds context to the original error. For example
```go
_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "IO_TEC_ReadFile", "read failed")
}
```
## Handling errors
You can easily handle different types of errors without identifying the source of the error.
```go
// Assume that access denied errors codes start with "AD_"
if errors.Match(err, regexp.MustCompile("^AD_.*$")) {
	// Handling access denied errors
}
```