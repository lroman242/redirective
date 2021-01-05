package logger_test

import (
	"strings"
	"testing"

	"github.com/lroman242/redirective/infrastructure/logger"
)

const expectedMessage = "Some Expected Message"

func TestStackTrace_Message(t *testing.T) {
	stackTrace := logger.NewStackTrace(expectedMessage, 1)

	if expectedMessage != stackTrace.Message() {
		t.Error("wrong message received")
	}
}

func TestStackTrace_Path(t *testing.T) {
	expectedPath := "/stacktrace_test.go:21"
	stackTrace := logger.NewStackTrace(expectedMessage, 1)

	if stackTrace.Path() != expectedPath {
		t.Errorf("wrong path received. expected %s but got %s", expectedPath, stackTrace.Path())
	}
}

func TestStackTrace_String(t *testing.T) {
	expectedPath := "/stacktrace_test.go:30"
	stackTrace := logger.NewStackTrace(expectedMessage, 1)

	if !strings.Contains(stackTrace.String(), expectedMessage) {
		t.Error("result string doesn't contains input message")
	}

	if !strings.Contains(stackTrace.String(), expectedPath) {
		t.Error("result string doesn't contains path")
	}
}
