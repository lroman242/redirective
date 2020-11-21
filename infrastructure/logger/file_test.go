package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	switch logger.(type) {
	case *fileLogger:
		if logger.(*fileLogger).logsDir != expectedLogsPath {
			t.Error("wrong logs dir path")
		}
		if _, err := os.Stat(logger.(*fileLogger).logsDir); os.IsNotExist(err) {
			t.Error("logs directory is not exists")
		}
		if _, err := os.Stat(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename); os.IsNotExist(err) {
			t.Error("log file is not exists")
		}
	default:
		t.Error("unexpected type")
	}
}

func TestFileLogger_Write(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	_, err := logger.Write([]byte(expectedStringLog))
	if err != nil {
		t.Errorf("unexpected error during write. error: %s", err)
	}

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if expectedStringLog != string(content) {
		t.Errorf("unexpected logs parsed. Expected %s but got %s", expectedStringLog, string(content))
	}
}

func TestFileLogger_Debug(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Debug(expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "debug") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "debug", string(content))
	}
}

func TestFileLogger_Info(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Info(expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "info") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "info", string(content))
	}
}

func TestFileLogger_Warn(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Warn(expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "warning") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "warning", string(content))
	}
}

func TestFileLogger_Error(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Error(expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "error") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "error", string(content))
	}
}

func TestFileLogger_Fatal(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Fatal(expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "fatal") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "fatal", string(content))
	}
}

func TestFileLogger_Panic(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Panic expected")
		}
	}()
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Panic(expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "panic") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "panic", string(content))
	}
}

func TestFileLogger_Debugf(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	expectedFormat := "log: %s"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Debugf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "debug") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "debug", string(content))
	}
}

func TestFileLogger_Infof(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	expectedFormat := "log: %s"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Infof(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "info") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "info", string(content))
	}
}

func TestFileLogger_Warnf(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	expectedFormat := "log: %s"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Warnf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "warning") {
		t.Errorf("unexpected logs parsed. expected substring %s in %s", "warning", string(content))
	}
}

func TestFileLogger_Errorf(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	expectedFormat := "log: %s"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Errorf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "error") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "fatal", string(content))
	}
}

func TestFileLogger_Fatalf(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	expectedFormat := "log: %s"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Fatalf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
	if !strings.Contains(string(content), "fatal") {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "fatal", string(content))
	}
}

func TestFileLogger_Printf(t *testing.T) {
	expectedStringLog := "some log message"
	expectedLogsPath := "test_logs"
	expectedFormat := "log: %s"
	logger := NewFileLogger(expectedLogsPath)
	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	logger.Printf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(logger.(*fileLogger).logsDir + "/" + logger.(*fileLogger).filename)
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
}
