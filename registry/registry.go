package registry

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/infrastructure/heartbeat"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/lroman242/redirective/infrastructure/storage"
	"github.com/lroman242/redirective/infrastructure/tracer"
	"github.com/lroman242/redirective/interface/api/controllers"
	ip "github.com/lroman242/redirective/interface/api/presenter"
	ir "github.com/lroman242/redirective/interface/api/repository"
	"github.com/lroman242/redirective/usecase/interactor"
	"github.com/lroman242/redirective/usecase/presenter"
	"github.com/lroman242/redirective/usecase/repository"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultScreenWidth = 1920
const defaultScreenHeight = 1080

type registry struct {
	conf      *config.AppConfig
	storage   storage.Storage
	logger    logger.Logger
	tracer    tracer.Tracer
	heartbeat heartbeat.HeartBeat
}

type Registry interface {
	NewHandler() http.Handler
	NewTraceController() controllers.TraceController
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
		Width:  defaultScreenWidth,
		Height: defaultScreenHeight,
	}, conf.ScreenshotsPath)

	//hb := heartbeat.NewProcessChecker(cmd.Process, fl)

	return &registry{
		conf:    conf,
		storage: mgdb,
		logger:  fl,
		tracer:  tr,
		//heartbeat: hb,
	}
}

func (r *registry) NewTraceController() controllers.TraceController {
	return controllers.NewTraceController(r.NewTraceInteractor(), r.conf.ScreenshotsPath, r.logger)
}

func (r *registry) NewTraceInteractor() interactor.TraceInteractor {
	return interactor.NewTraceInteractor(r.tracer, r.NewTracePresenter(), r.NewTracerRepository(), r.logger)
}

func (r *registry) NewTracerRepository() repository.TraceRepository {
	return ir.NewTraceRepository(r.storage)
}

func (r *registry) NewTracePresenter() presenter.TracePresenter {
	return ip.NewTracePresenter()
}

func (r *registry) NewHandler() http.Handler {
	router := httprouter.New()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// add routes
	router.GET("/api/find/:id", func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		r.logger.Printf("[%s] Find: %s", time.Now().Format(time.RFC3339), id)
		controllers.LoadTraceResults(writer, request, col, id)
	})
	router.GET("/api/screenshot/chrome", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		r.logger.Printf("[%s] Screenshot request: %s", time.Now().Format(time.RFC3339), request.URL.Query().Get("url"))
		controllers.ChromeScreenshot(writer, request, screenshotsStoragePath)
	})
	router.GET("/api/trace/chrome", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		r.logger.Printf("[%s] Trace request: %s", time.Now().Format(time.RFC3339), request.URL.Query().Get("url"))
		controllers.ChromeTrace(writer, request, screenshotsStoragePath, col)
	})

	// Serve static files from the ./assets directory
	// http(s)://api.redirective.net/screenshots/{filename.png}
	router.NotFound = http.FileServer(http.Dir("assets/"))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(router)

	// Insert the middleware
	handler = c.Handler(handler)

	return handler
}
