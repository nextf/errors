package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nextf/errors"
)

var ErrNotFoundPage = errors.ErrCode("NOT_FOUND", "Not found page")
var ErrNotFoundFile = errors.ErrCode("NOT_FOUND", "Not found file")
var ErrEndOfStream = errors.ErrCode("EOF", "End of stream")

func TestAs(t *testing.T) {
	e1 := errors.Wrap(ErrNotFoundPage, "Not found index.html")
	var coder interface{ Code() string }
	if errors.As(e1, &coder) {
		if coder.Code() != "NOT_FOUND" {
			t.Errorf("Expect %v, got %v", "code=NOT_FOUND", coder.Code())
		}
	} else {
		t.Errorf("Expect %v, got %v", "interface{ Code() string }", reflect.TypeOf(e1).String())
	}
}

func TestMatch(t *testing.T) {
	e1 := errors.Wrap(ErrNotFoundPage, "Not found index.html")
	if !errors.Match(e1, "NOT_FOUND") {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", "[NotMatch]")
	}
	if errors.Match(e1, "CODE0001") {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", "CODE0001")
	}
}

func TestFormat(t *testing.T) {
	if fmt.Sprintf("%q", ErrNotFoundPage) != "\"Not found page\"" {
		t.Errorf("Expect %q, got %q", "Not found page", ErrNotFoundPage)
	}
	if fmt.Sprintf("%s", ErrNotFoundPage) != "Not found page" {
		t.Errorf("Expect %s, got %s", "Not found page", ErrNotFoundPage)
	}
	if fmt.Sprintf("%v", ErrNotFoundPage) != "[NOT_FOUND] Not found page" {
		t.Errorf("Expect %v, got %v", "Not found page", ErrNotFoundPage)
	}
	if fmt.Sprintf("%+v", ErrNotFoundPage) != "[NOT_FOUND] Not found page" {
		t.Errorf("Expect %+v, got %+v", "Not found page", ErrNotFoundPage)
	}

	e1 := errors.Wrap(ErrNotFoundPage, "Not found index.html")
	withStack := fmt.Sprintf("%+v", e1)
	logPrefix := "[NOT_FOUND] Not found page\nNot found index.html\ngithub.com/nextf/errors.Wrap"
	if !strings.HasPrefix(withStack, logPrefix) {
		t.Error("Unprinted stack.")
	}
}

func TestWithStack(t *testing.T) {
	err := errors.ErrCode("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.WithStack(err)
	}
	counter := 0
	level := 0
	for {
		level++
		if _, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
			counter++
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	if counter != 1 {
		t.Errorf("Expect %d, got %d", 1, counter)
	}
	if level != 2 {
		t.Errorf("Expect %d, got %d", 1, level)
	}
}

func TestWrap(t *testing.T) {
	err := errors.ErrCode("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.Wrap(err, fmt.Sprintf("Level %d", i))
	}
	counter := 0
	level := 0
	for {
		level++
		if _, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
			counter++
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	if counter != 1 {
		t.Errorf("Expect %d, got %d", 1, counter)
	}
	if level != times+1 {
		t.Errorf("Expect %d, got %d", times+1, level)
	}
}

func TestWrapf(t *testing.T) {
	err := errors.ErrCode("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.Wrapf(err, "Level %d", i)
	}
	counter := 0
	level := 0
	for {
		level++
		if _, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
			counter++
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	if counter != 1 {
		t.Errorf("Expect %d, got %d", 1, counter)
	}
	if level != times+1 {
		t.Errorf("Expect %d, got %d", times+1, level)
	}
}
