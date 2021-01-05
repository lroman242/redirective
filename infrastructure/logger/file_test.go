package logger_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/lroman242/redirective/infrastructure/logger"
)

const (
	expectedStringLog = `some log message`
	expectedLogsPath  = `test_logs`
	expectedFormat    = `log: %s`
)

func TestNewFileLogger(t *testing.T) {
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	if _, err := os.Stat(log.(*logger.FileLogger).LogFilePath()); os.IsNotExist(err) {
		t.Error("log file is not exists")
	}
}

func TestFileLogger_Write(t *testing.T) {
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	_, err := log.Write([]byte(expectedStringLog))
	if err != nil {
		t.Errorf("unexpected error during write. error: %s", err)
	}

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if expectedStringLog != string(content) {
		t.Errorf("unexpected logs parsed. Expected %s but got %s", expectedStringLog, string(content))
	}
}

func TestFileLogger_Debug(t *testing.T) {
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Debug(expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Info(expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Warn(expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Error(expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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

//func TestFileLogger_Fatal(t *testing.T) {
//	log := logger.NewFileLogger(expectedLogsPath)
//	defer func() {
//		_ = os.RemoveAll(expectedLogsPath)
//	}()
//
//	log.Fatal(expectedStringLog)
//
//	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
//	if err != nil {
//		t.Errorf("unexpected error while reading log file: %s", err)
//	}
//
//	if !strings.Contains(string(content), expectedStringLog) {
//		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
//	}
//
//	if !strings.Contains(string(content), "fatal") {
//		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "fatal", string(content))
//	}
//}

func TestFileLogger_Panic(t *testing.T) {
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Panic expected")
		}
	}()

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Panic(expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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

func TestFileLogger_Debugf(t *testing.T) {
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Debugf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Infof(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Warnf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Errorf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
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

//func TestFileLogger_Fatalf(t *testing.T) {
//	log := logger.NewFileLogger(expectedLogsPath)
//	defer func() {
//		_ = os.RemoveAll(expectedLogsPath)
//	}()
//
//	log.Fatalf(expectedFormat, expectedStringLog)
//
//	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
//	if err != nil {
//		t.Errorf("unexpected error while reading log file: %s", err)
//	}
//
//	if !strings.Contains(string(content), expectedStringLog) {
//		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
//	}
//	if !strings.Contains(string(content), "error") {
//		t.Errorf("unexpected logs parsed. Expected substring %s in %s", "error", string(content))
//	}
//}

func TestFileLogger_Printf(t *testing.T) {
	log := logger.NewFileLogger(expectedLogsPath)

	defer func() {
		_ = os.RemoveAll(expectedLogsPath)
	}()

	log.Printf(expectedFormat, expectedStringLog)

	content, err := ioutil.ReadFile(log.(*logger.FileLogger).LogFilePath())
	if err != nil {
		t.Errorf("unexpected error while reading log file: %s", err)
	}

	if !strings.Contains(string(content), expectedStringLog) {
		t.Errorf("unexpected logs parsed. Expected substring %s in %s", expectedStringLog, string(content))
	}
}
