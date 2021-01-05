package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const traceLevel = 2

// FileLogger type implement Logger interface and write logs to local file.
type FileLogger struct {
	logsDir  string
	lock     sync.Mutex
	filename string // should be set to the actual filename
	fp       *os.File
}

// NewFileLogger function build new instance of fileLogger which implements Logger interface.
func NewFileLogger(logsDirPath string) Logger {
	if _, err := os.Stat(logsDirPath); os.IsNotExist(err) {
		// logs directory does not exist
		err = os.Mkdir(logsDirPath, 0o750)
		if err != nil {
			panic(err)
		}
	}

	l := &FileLogger{
		logsDir: logsDirPath,
	}
	l.filename = l.fileNameForNow()

	err := l.rotate()
	if err != nil {
		panic(err)
	}

	return l
}

// LogFilePath function return path to current log file.
func (l *FileLogger) LogFilePath() string {
	return l.logsDir + "/" + l.filename
}

func (l *FileLogger) fileNameForNow() string {
	return fmt.Sprintf("%s_redirective.log", time.Now().Format("2006-01-02"))
}

func (l *FileLogger) prefix(level string) string {
	return fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), level)
}

// Write function satisfies the io.Writer interface.
func (l *FileLogger) Write(output []byte) (int, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.filename != l.fileNameForNow() {
		err := l.rotate()
		if err != nil {
			panic(err)
		}
	}

	return l.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (l *FileLogger) rotate() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Close existing file if open
	if l.fp != nil {
		err := l.fp.Close()
		l.fp = nil

		if err != nil {
			return err
		}
	}

	l.filename = l.fileNameForNow()

	// Rename dest file if it already exists
	_, err := os.Stat(l.LogFilePath())
	if err == nil {
		err = os.Rename(l.LogFilePath(), l.LogFilePath()+"."+time.Now().Format(time.RFC3339))
		if err != nil {
			return err
		}
	}

	// Create a file.
	l.fp, err = os.Create(l.LogFilePath())
	if err != nil {
		panic(err)
	}

	return nil
}

// Debugf write formatted message with DEBUG level.
func (l *FileLogger) Debugf(format string, data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("debug") + NewStackTrace(fmt.Sprintf(format, data...), traceLevel).String() + "\n"))
}

// Infof write formatted message with INFO level.
func (l *FileLogger) Infof(format string, data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("info") + NewStackTrace(fmt.Sprintf(format, data...), traceLevel).String() + "\n"))
}

// Printf write formatted message.
func (l *FileLogger) Printf(format string, data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("-") + NewStackTrace(fmt.Sprintf(format, data...), traceLevel).String() + "\n"))
}

// Warnf write formatted message with WARNING level.
func (l *FileLogger) Warnf(format string, data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("warning") + NewStackTrace(fmt.Sprintf(format, data...), traceLevel).String() + "\n"))
}

// Errorf write formatted message with ERROR level.
func (l *FileLogger) Errorf(format string, data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("error") + NewStackTrace(fmt.Sprintf(format, data...), traceLevel).String() + "\n"))
}

// Fatalf write formatted message with ERROR level and exit.
func (l *FileLogger) Fatalf(format string, data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("error") + NewStackTrace(fmt.Sprintf(format, data...), traceLevel).String() + "\n"))

	os.Exit(1)
}

// Debug write message with DEBUG level.
func (l *FileLogger) Debug(data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("debug") + NewStackTrace(fmt.Sprint(data...), traceLevel).String() + "\n"))
}

// Info write message with INFO level.
func (l *FileLogger) Info(data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("info") + NewStackTrace(fmt.Sprint(data...), traceLevel).String() + "\n"))
}

// Warn write message with WARNING level.
func (l *FileLogger) Warn(data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("warning") + NewStackTrace(fmt.Sprint(data...), traceLevel).String() + "\n"))
}

// Error write message with ERROR level.
func (l *FileLogger) Error(data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("error") + NewStackTrace(fmt.Sprint(data...), traceLevel).String() + "\n"))
}

// Fatal write message with ERROR level and exit.
func (l *FileLogger) Fatal(data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("fatal") + NewStackTrace(fmt.Sprint(data...), traceLevel).String() + "\n"))

	os.Exit(1)
}

// Panic write message with ERROR level and throw panic exception.
func (l *FileLogger) Panic(data ...interface{}) {
	_, _ = l.Write([]byte(l.prefix("error") + NewStackTrace(fmt.Sprint(data...), traceLevel).String() + "\n"))
	panic(data)
}
