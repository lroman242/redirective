package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// StackTrace type describe error message with simple stack trace.
type StackTrace struct {
	msg  string
	path string
}

// NewStackTrace function constructs a new StackTrace struct by using given panic
// message, absolute path of the caller file and the line number.
func NewStackTrace(msg string, lvl int) *StackTrace {
	_, file, line, _ := runtime.Caller(lvl)
	p, _ := os.Getwd()

	return &StackTrace{
		msg:  msg,
		path: fmt.Sprintf("%s:%d", strings.TrimPrefix(file, p), line),
	}
}

// Message function provide access to message field.
func (s *StackTrace) Message() string {
	return s.msg
}

// Path function provide access to trace path.
func (s *StackTrace) Path() string {
	return s.path
}

// String function convert instance to string. Satisfy Stringer interface.
func (s *StackTrace) String() string {
	return fmt.Sprintf("%s: %s", s.path, s.Message())
}
