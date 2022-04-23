package errors_test

import (
	"testing"

	"github.com/nextf/errors"
)

const (
	errNoCode        = errors.ConstError("No code")
	errInvalidCode   = errors.ConstError("[ERR 001] This error has a invalid code")
	errNormalize     = errors.ConstError("[ERR001] This error has a code")
	errAbnormality   = errors.ConstError("\t[ERR001]   This error has a code  ")
	errHasSeparator1 = errors.ConstError("[ERR-001] This error has a code")
	errHasSeparator2 = errors.ConstError("[ERR_001] This error has a code")
)

func TestConstError(t *testing.T) {
	var err error = errNoCode
	code, ok := errors.GetCode(err)
	if !ok {
		t.Errorf("No code in ConstError")
	}
	if code != "" {
		t.Errorf("Get unexpected code")
	}
	if err.Error() != "No code" {
		t.Errorf("Unexpected error message")
	}
	if !errors.Match(err, "") {
		t.Errorf("Get unexpected code")
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
	if !errors.Match(err, "") {
		t.Errorf("Get unexpected code")
	}

	err = errNormalize
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
	if !errors.Match(err, "ERR001") {
		t.Errorf("Get unexpected code")
	}

	err = errAbnormality
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
	if !errors.Match(err, "ERR001") {
		t.Errorf("Get unexpected code")
	}

	err = errHasSeparator1
	code, ok = errors.GetCode(err)
	if !ok {
		t.Errorf("No code in ConstError")
	}
	if code != "ERR-001" {
		t.Errorf("Get unexpected code")
	}
	if err.Error() != "This error has a code" {
		t.Errorf("Unexpected error message")
	}
	if !errors.Match(err, "ERR-001") {
		t.Errorf("Get unexpected code")
	}

	err = errHasSeparator2
	code, ok = errors.GetCode(err)
	if !ok {
		t.Errorf("No code in ConstError")
	}
	if code != "ERR_001" {
		t.Errorf("Get unexpected code")
	}
	if err.Error() != "This error has a code" {
		t.Errorf("Unexpected error message")
	}
	if !errors.Match(err, "ERR_001") {
		t.Errorf("Get unexpected code")
	}
}
