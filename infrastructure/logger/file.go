package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type fileLogger struct {
	logsDir  string
	lock     sync.Mutex
	filename string // should be set to the actual filename
	fp       *os.File
}

func NewFileLogger(logsDirPath string) *fileLogger {
	if _, err := os.Stat(logsDirPath); os.IsNotExist(err) {
		// logs directory does not exist
		err = os.Mkdir(logsDirPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	l := &fileLogger{
		logsDir: logsDirPath,
	}
	l.filename = l.fileNameForNow()

	err := l.Rotate()
	if err != nil {
		panic(err)
	}

	return l
}

func (l *fileLogger) logFilePath() string {
	return l.logsDir + "/" + l.filename
}

func (l *fileLogger) fileNameForNow() string {
	return fmt.Sprintf("%s_redirective.log", time.Now().Format("2006-01-02"))
}

func (l *fileLogger) prefix(level string) string {
	return fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), level)
}

// Write satisfies the io.Writer interface.
func (l *fileLogger) Write(output []byte) (int, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.filename != l.fileNameForNow() {
		err := l.Rotate()
		if err != nil {
			panic(err)
		}
	}

	return l.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (l *fileLogger) Rotate() (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Close existing file if open
	if l.fp != nil {
		err = l.fp.Close()
		l.fp = nil
		if err != nil {
			return
		}
	}
	l.filename = l.fileNameForNow()

	// Rename dest file if it already exists
	_, err = os.Stat(l.logFilePath())
	if err == nil {
		err = os.Rename(l.logFilePath(), l.logFilePath()+"."+time.Now().Format(time.RFC3339))
		if err != nil {
			return
		}
	}

	// Create a file.
	l.fp, err = os.Create(l.logFilePath())
	if err != nil {
		panic(err)
	}

	return
}

func (l *fileLogger) Debugf(format string, data ...interface{}) {
	l.Write([]byte(l.prefix("debug") +  NewStackTrace(fmt.Sprintf(format+"\n", data...)).String() ))
}

func (l *fileLogger) Infof(format string, data ...interface{}) {
	l.Write([]byte(l.prefix("info") + NewStackTrace(fmt.Sprintf(format+"\n", data...)).String()))
}

func (l *fileLogger) Printf(format string, data ...interface{}) {
	l.Write([]byte(l.prefix("-") + NewStackTrace(fmt.Sprintf(format+"\n", data...)).String()))
}

func (l *fileLogger) Warnf(format string, data ...interface{}) {
	l.Write([]byte(l.prefix("warn") + NewStackTrace(fmt.Sprintf(format+"\n", data...)).String()))
}

func (l *fileLogger) Errorf(format string, data ...interface{}) {
	l.Write([]byte(l.prefix("error") + NewStackTrace(fmt.Sprintf(format+"\n", data...)).String()))
}

func (l *fileLogger) Fatalf(format string, data ...interface{}) {
	l.Write([]byte(l.prefix("fatal") + NewStackTrace(fmt.Sprintf(format+"\n", data...)).String()))
}

func (l *fileLogger) Debug(data ...interface{}) {
	l.Write([]byte(l.prefix("debug") + NewStackTrace(fmt.Sprint(data...)).String()))
}

func (l *fileLogger) Info(data ...interface{}) {
	l.Write([]byte(l.prefix("info") + NewStackTrace(fmt.Sprint(data...)).String()))
}

func (l *fileLogger) Warn(data ...interface{}) {
	l.Write([]byte(l.prefix("warn") + NewStackTrace(fmt.Sprint(data...)).String()))
}

func (l *fileLogger) Error(data ...interface{}) {
	l.Write([]byte(l.prefix("error") + NewStackTrace(fmt.Sprint(data...)).String()))
}

func (l *fileLogger) Fatal(data ...interface{}) {
	l.Write([]byte(l.prefix("fatal") + NewStackTrace(fmt.Sprint(data...)).String()))
}

func (l *fileLogger) Panic(data ...interface{}) {
	l.Write([]byte(l.prefix("panic") + NewStackTrace(fmt.Sprint(data...)).String()))
}
