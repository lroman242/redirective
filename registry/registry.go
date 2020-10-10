package registry

import (
	"fmt"
	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/infrastructure/heartbeat"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/lroman242/redirective/infrastructure/storage"
	"github.com/lroman242/redirective/infrastructure/tracer"
	"log"
	"os"
)

type registry struct {
	storage storage.Storage
	logger logger.Logger
	tracer tracer.Tracer
	heartbeat heartbeat.HeartBeat
}

//TODO:
type Registry interface {
	//NewAppController() controller.AppController
}

func NewRegistry(conf *config.AppConfig) Registry {
	if _, err := os.Stat(conf.LogFilePath); os.IsNotExist(err) {
		// logs directory does not exist
		err = os.Mkdir(conf.LogFilePath, 0755)
		if err != nil {
			panic(fmt.Sprintf("Logs dirrectory (%s) is not exists and couldn't be created. Error: %s", conf.LogFilePath, err))
		}
	}

	fl := logger.NewFileLogger(conf.LogFilePath)
	log.SetOutput(fl)

	mgdb, err := storage.NewMongoDB(conf.Storage)
	if err != nil {
		panic(fmt.Sprintf("Mongo DB storage couldn't be initialized. Error: %s", err))
	}

	if _, err := os.Stat(conf.ScreenshotsPath); os.IsNotExist(err) {
		// logs directory does not exist
		err = os.Mkdir(conf.ScreenshotsPath, 0755)
		if err != nil {
			panic(fmt.Sprintf("Screenshots dirrectory (%s) is not exists and couldn't be created. Error: %s", conf.ScreenshotsPath, err))
		}
	}

	tr := tracer.NewChromeTracer(&tracer.ScreenSize{
		Width:  1920,
		Height: 1080,
	}, conf.ScreenshotsPath)

	//hb := heartbeat.NewProcessChecker(cmd.Process, fl)

	return &registry{
		storage:   mgdb,
		logger:    fl,
		tracer:    tr,
		//heartbeat: hb,
	}
}