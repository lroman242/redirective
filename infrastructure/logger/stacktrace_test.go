package logger

import (
	"strings"
	"testing"
)

func TestNewStackTrace(t *testing.T) {
	expectedMessage := "Some Expected Message"
	expectedPath := "/stacktrace_test.go:11"
	stackTrace := NewStackTrace(expectedMessage, 1)

	if stackTrace.msg != expectedMessage {
		t.Error("wrong message received")
	}
	if stackTrace.path != expectedPath {
		t.Error("wrong path received")
	}
}

func TestStackTrace_Message(t *testing.T) {
	expectedMessage := "Some Expected Message"
	stackTrace := NewStackTrace(expectedMessage, 1)

	if expectedMessage != stackTrace.Message() {
		t.Error("wrong message received")
	}
}

func TestStackTrace_Path(t *testing.T) {
	expectedMessage := "Some Expected Message"
	expectedPath := "/stacktrace_test.go:33"
	stackTrace := NewStackTrace(expectedMessage, 1)

	if stackTrace.Path() != expectedPath {
		t.Errorf("wrong path received. expected %s but got %s", expectedPath, stackTrace.Path())
	}
}

func TestStackTrace_String(t *testing.T) {
	expectedMessage := "Some Expected Message"
	expectedPath := "/stacktrace_test.go:43"
	stackTrace := NewStackTrace(expectedMessage, 1)

	if !strings.Contains(stackTrace.String(), expectedMessage) {
		t.Error("result string doesn't contains input message")
	}
	if !strings.Contains(stackTrace.String(), expectedPath) {
		t.Error("result string doesn't contains path")
	}
}
