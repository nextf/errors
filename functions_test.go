package errors_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/nextf/errors"
	"github.com/nextf/errors/stack"
)

var ErrNotFoundPage = errors.ErrCode("NOT_FOUND", "Not found page")
var ErrEndOfStream = errors.ErrCode("EOF", "End of stream")

func TestErrCodeAs(t *testing.T) {
	e1 := errors.TraceMessage(ErrNotFoundPage, "Not found index.html")
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
	e1 := errors.TraceMessage(ErrNotFoundPage, "Not found index.html")
	if !errors.Match(e1, "NOT_FOUND") {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", "[NotMatch]")
	}
	if errors.Match(e1, "CODE0001") {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", "CODE0001")
	}
	if errors.Match(e1, []byte("NOT_FOUND")) {
		t.Errorf("Expect %v, got %v", "code=NOT_FOUND", []byte("NOT_FOUND"))
	}
	if !errors.Match(e1, regexp.MustCompile("^NOT.*$")) {
		t.Errorf("Expect %v, got %v", "code~=^NOT.*$", "[NotMatch]")
	}
	if errors.Match(e1, regexp.MustCompile("^YES.*$")) {
		t.Errorf("Expect %v, got %v", "code~=^NOT.*$", "[NotMatch]")
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

	e1 := errors.TraceMessage(ErrNotFoundPage, "Not found index.html")
	withStack := fmt.Sprintf("%+v", e1)
	logPrefix := "Not found index.html\nCaused by: @callstack\n\x20\x20\x20\x20github.com/nextf/errors_test.TestErrCodeFormat(functions_test.go:63)"
	if !strings.HasPrefix(withStack, logPrefix) {
		t.Error("Unprinted stack.")
	}
	e2 := errors.WithErrCode(errors.New("Not found file"), "NOT_FOUND", "Not found index.html")
	withCode := fmt.Sprintf("%+v", e2)
	logText := "[NOT_FOUND] Not found index.html\nCaused by: Not found file"
	if withCode != logText {
		t.Errorf("Expected message was not obtained\n%s", withCode)
	}
}

func TestWithStackNodup(t *testing.T) {
	err := errors.ErrCode("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.TraceNodup(err)
	}
	counter := 0
	level := 0
	for {
		if _, ok := err.(interface{ StackTrace() []stack.Frame }); ok {
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

func TestTraceMessage(t *testing.T) {
	err := errors.ErrCode("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.TraceMessage(err, fmt.Sprintf("Level %d", i))
	}
	counter := 0
	level := 0
	for {
		if _, ok := err.(interface{ StackTrace() []stack.Frame }); ok {
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

func TestTraceMessagef(t *testing.T) {
	err := errors.ErrCode("ROOT", "Level 0")
	times := 10
	for i := 1; i < times; i++ {
		err = errors.TraceMessagef(err, "Level %d", i)
	}
	counter := 0
	level := 0
	for {
		if _, ok := err.(interface{ StackTrace() []stack.Frame }); ok {
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

func TestErrCode(t *testing.T) {
	err := errors.ErrCode("NEW_IN_FUNC", "error for unittest")
	if errors.HasErrorStack(err) {
		t.Errorf("It was expected that there would be no StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != "error for unittest" {
		t.Errorf("Expect `%s`, got `%s`", "error for unittest", err.Error())
	}
}

func TestTraceableErrCode(t *testing.T) {
	err := errors.TraceableErrCode("NEW_IN_FUNC", "error for unittest")
	if !errors.HasErrorStack(err) {
		t.Errorf("It was expected that there would has StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != "error for unittest" {
		t.Errorf("Expect `%s`, got `%s`", "error for unittest", err.Error())
	}
}

func TestErrCodef(t *testing.T) {
	num := rand.Int()
	err := errors.ErrCodef("NEW_IN_FUNC", "Wrong number [%d]", num)
	if errors.HasErrorStack(err) {
		t.Errorf("It was expected that there would be no StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != fmt.Sprintf("Wrong number [%d]", num) {
		t.Errorf("Expect `%s`, got `%s`", fmt.Sprintf("Wrong number [%d]", num), err.Error())
	}
}

func TestTraceableErrCodef(t *testing.T) {
	num := rand.Int()
	err := errors.TraceableErrCodef("NEW_IN_FUNC", "Wrong number [%d]", num)
	if !errors.HasErrorStack(err) {
		t.Errorf("It was expected that there would has StackTrace in the `err`, but it wasn't.")
	}
	if code, ok := errors.GetCode(err); !ok || code != "NEW_IN_FUNC" {
		t.Errorf("Expect code=%s, got [%s]", "NEW_IN_FUNC", code)
	}
	if err.Error() != fmt.Sprintf("Wrong number [%d]", num) {
		t.Errorf("Expect `%s`, got `%s`", fmt.Sprintf("Wrong number [%d]", num), err.Error())
	}
}

func TestWrapNodup(t *testing.T) {
	err := errors.New("[ROOT_ERROR] Root level")
	times := 10
	for i := 1; i <= times; i++ {
		err = errors.WrapNodup(err, fmt.Sprintf("L%d", i), fmt.Sprintf("Level %d", i))
	}
	stackCounter := 0
	codeCounter := 0
	level := 0
	errTmp := err
	for {
		if _, ok := errTmp.(interface{ StackTrace() []stack.Frame }); ok {
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
	if codeCounter != times+1 {
		t.Errorf("Expect %d, got %d", times+1, codeCounter)
	}
	if level != times+1 {
		t.Errorf("Expect %d, got %d", times+1, level)
	}
}

func TestWrap(t *testing.T) {
	err := errors.New("[ROOT_ERROR] Root level")
	times := 10
	for i := 1; i <= times; i++ {
		err = errors.Wrap(err, fmt.Sprintf("L%d", i), fmt.Sprintf("Level %d", i))
	}
	stackCounter := 0
	codeCounter := 0
	level := 0
	errTmp := err
	for {
		if _, ok := errTmp.(interface{ StackTrace() []stack.Frame }); ok {
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
	if stackCounter != times {
		t.Errorf("Expect %d, got %d", times, stackCounter)
	}
	if codeCounter != times+1 {
		t.Errorf("Expect %d, got %d", times+1, codeCounter)
	}
	if level != times*2 {
		t.Errorf("Expect %d, got %d", times*2, level)
	}
}
