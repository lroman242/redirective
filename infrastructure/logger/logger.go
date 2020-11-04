package logger

//go:generate mockgen -package=mocks -destination=mocks/mock_logger.go -source=infrastructure/logger/logger.go Logger

// Logger interface describe logger instance used in application to process logs
type Logger interface {
	Write([]byte) (int, error)
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}
