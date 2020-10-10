package heartbeat

import (
	"context"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/opentracing/opentracing-go/log"
	"os"
	"syscall"
	"time"
)

// processChecker is implementation of HeartBeat which check
// if all processes are alive using os.Process instance
type processChecker struct {
	chromeProcess *os.Process
	log logger.Logger
}

// NewProcessChecker function will create new instance of processChecker
func NewProcessChecker(process *os.Process, log logger.Logger) HeartBeat {
	return &processChecker{
		chromeProcess: process,
		log: log,
	}
}

// Check func will check if all services
// required for application are "alive"
func (pc *processChecker) Check() error {
	err := pc.checkChromeProcess()
	if err != nil {
		return err
	}

	// TODO: check other services required for app
	return nil
}

// Monitor function will check if application is alive (using Check method)
// each "duration" piece of time
func (pc *processChecker) Monitor(ctx context.Context, duration time.Duration) {
	go func() {
		for {
			select {
			case <- ctx.Done():
				return // break infinity loop
			case <- time.After(duration):
				err := pc.Check()
				if err != nil {
					log.Error(err)
					panic(err)
				}
			}
		}
	}()
}

// checkChromeProcess will do a check if google chrome process is alive
// and return error if something is wrong
func (pc *processChecker)checkChromeProcess() error {
	err := pc.chromeProcess.Signal(syscall.Signal(0))
	if err != nil {
		errno, ok := err.(syscall.Errno)
		if !ok {
			return err
		}
		switch errno {
		case syscall.ESRCH:
			return err
		case syscall.EPERM:
			return nil
		}

		return err
	}

	return nil
}