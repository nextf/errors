package errors_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/nextf/errors"
)

var ErrNotFoundPage = errors.New("NOT_FOUND", "Not found page")
var ErrEndOfStream = errors.New("EOF", "End of stream")

func TestErrCodeAs(t *testing.T) {
	e1 := errors.WithMessage(ErrNotFoundPage, "Not found index.html")
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
	e1 := errors.WithMessage(ErrNotFoundPage, "Not found index.html")
	if !errors.Match(e1, "NOT_FOUND") {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", "[NotMatch]")
	}
	if errors.Match(e1, "CODE0001") {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", "CODE0001")
	}
}

func TestErrCodeFormat(t *testing.T) {
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

	e1 := errors.WithMessage(ErrNotFoundPage, "Not found index.html")
	withStack := fmt.Sprintf("%+v", e1)
	logPrefix := "[NOT_FOUND] Not found page\nNot found index.html\ngithub.com/nextf/errors_test.TestErrCodeFormat\n"
	if !strings.HasPrefix(withStack, logPrefix) {
		t.Error("Unprinted stack.")
	}
	e2 := errors.WithCode(errors.Message("Not found file"), "NOT_FOUND", "Not found index.html")
	withCode := fmt.Sprintf("%+v", e2)
	logText := "[NOT_FOUND] Not found index.html\nCaused by: Not found file"
	if withCode != logText {
		t.Errorf("Expected message was not obtained\n%s", withCode)
	}
}

func TestWithStack(t *testing.T) {
	err := errors.New("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.WithStack(err)
	}
	counter := 0
	level := 0
	for {
		if _, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
			counter++
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
		level++
	}
	if counter != 1 {
		t.Errorf("Expect %d, got %d", 1, counter)
	}
	if level != 1 {
		t.Errorf("Expect %d, got %d", 1, level)
	}
}

func TestWithMessage(t *testing.T) {
	err := errors.New("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.WithMessage(err, fmt.Sprintf("Level %d", i))
	}
	counter := 0
	level := 0
	for {
		if _, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
			counter++
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
		level++
	}
	if counter != 1 {
		t.Errorf("Expect %d, got %d", 1, counter)
	}
	if level != times {
		t.Errorf("Expect %d, got %d", times, level)
	}
}

func TestWithMessagef(t *testing.T) {
	err := errors.New("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.WithMessagef(err, "Level %d", i)
	}
	counter := 0
	level := 0
	for {
		if _, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
			counter++
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
		level++
	}
	if counter != 1 {
		t.Errorf("Expect %d, got %d", 1, counter)
	}
	if level != times {
		t.Errorf("Expect %d, got %d", times, level)
	}
}

func TestNew(t *testing.T) {
	err := errors.New("NEW_IN_FUNC", "error for unittest")
	if errors.HasStackTrace(err) {
		t.Errorf("It was expected that there would be no StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != "error for unittest" {
		t.Errorf("Expect `%s`, got `%s`", "error for unittest", err.Error())
	}
}

func TestNewWithStack(t *testing.T) {
	err := errors.NewWithStack("NEW_IN_FUNC", "error for unittest")
	if !errors.HasStackTrace(err) {
		t.Errorf("It was expected that there would has StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != "error for unittest" {
		t.Errorf("Expect `%s`, got `%s`", "error for unittest", err.Error())
	}
}

func TestNewf(t *testing.T) {
	num := rand.Int()
	err := errors.Newf("NEW_IN_FUNC", "Wrong number [%d]", num)
	if errors.HasStackTrace(err) {
		t.Errorf("It was expected that there would be no StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != fmt.Sprintf("Wrong number [%d]", num) {
		t.Errorf("Expect `%s`, got `%s`", fmt.Sprintf("Wrong number [%d]", num), err.Error())
	}
}

func TestNewWithStackf(t *testing.T) {
	num := rand.Int()
	err := errors.NewWithStackf("NEW_IN_FUNC", "Wrong number [%d]", num)
	if !errors.HasStackTrace(err) {
		t.Errorf("It was expected that there would has StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != fmt.Sprintf("Wrong number [%d]", num) {
		t.Errorf("Expect `%s`, got `%s`", fmt.Sprintf("Wrong number [%d]", num), err.Error())
	}
}

func TestWrap_1(t *testing.T) {
	err := errors.New("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.Wrap(err, fmt.Sprintf("L%d", i), fmt.Sprintf("Level %d", i))
	}
	stackCounter := 0
	codeCounter := 0
	level := 0
	errTmp := err
	for {
		if _, ok := errTmp.(interface{ StackTrace() errors.StackTrace }); ok {
			stackCounter++
		}
		if _, ok := errTmp.(interface{ Code() string }); ok {
			codeCounter++
		}
		if errTmp = errors.Unwrap(errTmp); errTmp == nil {
			break
		}
		level++
	}
	if stackCounter != 1 {
		t.Errorf("Expect %d, got %d", 1, stackCounter)
	}
	if codeCounter != 1 {
		t.Errorf("Expect %d, got %d", 1, codeCounter)
	}
	if level != 1 {
		t.Errorf("Expect %d, got %d", 1, level)
	}
}

func TestWrap_2(t *testing.T) {
	err := errors.Message("Root level")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.Wrap(err, fmt.Sprintf("L%d", i), fmt.Sprintf("Level %d", i))
	}
	stackCounter := 0
	codeCounter := 0
	level := 0
	errTmp := err
	for {
		if _, ok := errTmp.(interface{ StackTrace() errors.StackTrace }); ok {
			stackCounter++
		}
		if _, ok := errTmp.(interface{ Code() string }); ok {
			codeCounter++
		}
		if errTmp = errors.Unwrap(errTmp); errTmp == nil {
			break
		}
		level++
	}
	if stackCounter != 1 {
		t.Errorf("Expect %d, got %d", 1, stackCounter)
	}
	if codeCounter != 1 {
		t.Errorf("Expect %d, got %d", 1, codeCounter)
	}
	if level != 2 {
		t.Errorf("Expect %d, got %d", 2, level)
	}
}
