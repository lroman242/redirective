// Package registry contains functions to initialize application services (web server, logger, etc)
package registry

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/config"
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
)

const (
	defaultScreenWidth  = 1920
	defaultScreenHeight = 1080
)

type registry struct {
	conf    *config.AppConfig
	storage storage.Storage
	logger  logger.Logger
	tracer  tracer.Tracer
}

// Registry interface describe struct which should build main application parts.
type Registry interface {
	NewHandler() http.Handler
	NewTraceController() controllers.TraceController
}

// NewRegistry function initialize new Registry instance.
func NewRegistry(conf *config.AppConfig) Registry {
	if _, err := os.Stat(conf.LogsPath); os.IsNotExist(err) {
		// logs directory does not exist
		err = os.Mkdir(conf.LogsPath, 0o750)
		if err != nil {
			panic(fmt.Sprintf("Logs dirrectory (%s) is not exists and couldn't be created. Error: %s", conf.LogsPath, err))
		}
	}

	fl := logger.NewFileLogger(conf.LogsPath)
	log.SetOutput(fl)

	mgdb, err := storage.NewMongoDB(conf.Storage)
	if err != nil {
		panic(fmt.Sprintf("Mongo DB storage couldn't be initialized. Error: %s", err))
	}

	if _, err := os.Stat(conf.ScreenshotsPath); os.IsNotExist(err) {
		// logs directory does not exist
		err = os.Mkdir(conf.ScreenshotsPath, 0o750)
		if err != nil {
			panic(fmt.Sprintf("Screenshots dirrectory (%s) is not exists and couldn't be created. Error: %s", conf.ScreenshotsPath, err))
		}
	}

	tr := tracer.NewChromeTracer(&tracer.ScreenSize{
		Width:  defaultScreenWidth,
		Height: defaultScreenHeight,
	}, conf.ScreenshotsPath)

	return &registry{
		conf:    conf,
		storage: mgdb,
		logger:  fl,
		tracer:  tr,
	}
}

// NewTraceController function will build controllers.TraceController instance.
func (r *registry) NewTraceController() controllers.TraceController {
	return controllers.NewTraceController(r.NewTraceInteractor(), r.conf.ScreenshotsPath, r.logger)
}

// NewTraceInteractor function will build interactor.TraceInteractor instance.
func (r *registry) NewTraceInteractor() interactor.TraceInteractor {
	return interactor.NewTraceInteractor(r.tracer, r.NewTracePresenter(), r.NewTracerRepository(), r.logger)
}

// NewTracerRepository function will build new repository.TraceRepository instance.
func (r *registry) NewTracerRepository() repository.TraceRepository {
	return ir.NewTraceRepository(r.storage)
}

// NewTracePresenter function will build new presenter.TracePresenter instance.
func (r *registry) NewTracePresenter() presenter.TracePresenter {
	var protocol string
	if r.conf.HTTPServer.HTTPS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	return ip.NewTracePresenter(r.conf.AppDomain, protocol)
}

// NewHandler will build new http.Handler with applied routes.
func (r *registry) NewHandler() http.Handler {
	router := httprouter.New()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	controller := r.NewTraceController()

	// add routes
	router.GET("/api/find/:id", controller.FindTraceResults)
	router.GET("/api/screenshot", controller.Screenshot)
	router.GET("/api/trace", controller.TraceURL)

	// Serve static files from the ./assets/screenshots directory
	// http(s)://api.redirective.net/screenshots/{filename.png}
	router.NotFound = http.StripPrefix("/screenshots", http.FileServer(http.Dir(r.conf.ScreenshotsPath)))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(router)

	// Insert the middleware
	handler = c.Handler(handler)

	return handler
}
