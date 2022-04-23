package errors_test

import (
	"testing"

	"github.com/nextf/errors"
)

const (
	errNoCode      = errors.ConstError("No code")
	errInvalidCode = errors.ConstError("[ERR 001] This error has a invalid code")
	errHasCode     = errors.ConstError("[ERR001] This error has a code ")
)

func TestConst(t *testing.T) {
	var err error = errNoCode
	code, ok := errors.GetCode(err)
	if !ok {
		t.Errorf("No code in ConstError")
	}
	if code != "" {
		t.Errorf("Get code by accident")
	}
	if err.Error() != "No code" {
		t.Errorf("Unexpected error message")
	}

	err = errInvalidCode
	code, ok = errors.GetCode(err)
	if !ok {
		t.Errorf("No code in ConstError")
	}
	if code != "" {
		t.Errorf("Get unexpected code")
	}
	if err.Error() != "[ERR 001] This error has a invalid code" {
		t.Errorf("Unexpected error message")
	}

	err = errHasCode
	code, ok = errors.GetCode(err)
	if !ok {
		t.Errorf("No code in ConstError")
	}
	if code != "ERR001" {
		t.Errorf("Get unexpected code")
	}
	if err.Error() != "This error has a code" {
		t.Errorf("Unexpected error message")
	}
}
