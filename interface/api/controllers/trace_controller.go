// Package controllers contains classes to handle http requests and build response
package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/lroman242/redirective/usecase/interactor"
)

const (
	defaultScreenWidth  = 1920
	defaultScreenHeight = 1080
)

// TraceController interface represent functions required to handle application requests.
type TraceController interface {
	TraceURL(http.ResponseWriter, *http.Request, httprouter.Params)
	Screenshot(http.ResponseWriter, *http.Request, httprouter.Params)
	FindTraceResults(http.ResponseWriter, *http.Request, httprouter.Params)
}

// traceController implement TraceController interface.
type traceController struct {
	traceInteractor interactor.TraceInteractor
	assetsPath      string
	log             logger.Logger
}

// NewTraceController will build TraceController instance.
func NewTraceController(ti interactor.TraceInteractor, assetsPath string, log logger.Logger) TraceController {
	return &traceController{ti, assetsPath, log}
}

// TraceURL function will trace redirects for provided URL.
func (tc *traceController) TraceURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}).Failed(w)

		return
	}

	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}).Failed(w)

		return
	}

	results, err := tc.traceInteractor.Trace(targetURL, tc.assetsPath)
	if err != nil {
		tc.log.Error(err)
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("an error occurred. error: %s", err),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}).Failed(w)

		return
	}

	(&Response{
		Status:     true,
		Message:    "url traced",
		StatusCode: http.StatusOK,
		Data:       results,
	}).Success(w)
}

// Screenshot function will retrieve screenshot from provided URL.
func (tc *traceController) Screenshot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}).Failed(w)

		return
	}

	width, height := parseSizeFromRequest(r)

	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}).Failed(w)

		return
	}

	screenshotURL, err := tc.traceInteractor.Screenshot(targetURL, width, height, tc.assetsPath)
	if err != nil {
		tc.log.Error(err)
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("an error occurred. error: %s", err),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}).Failed(w)

		return
	}

	(&Response{
		Status:     true,
		Message:    "url traced",
		StatusCode: http.StatusOK,
		Data:       screenshotURL,
	}).Success(w)
}

// FindTraceResults will find tracer results by provided id.
func (tc *traceController) FindTraceResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	results, err := tc.traceInteractor.FindTrace(id)
	if err != nil {
		tc.log.Error(err)
		(&Response{
			Status:     false,
			Message:    "trace results not found", //fmt.Sprintf("an error occurred. error: %s", err),
			StatusCode: http.StatusNotFound,
			Data:       nil,
		}).Failed(w)

		return
	}

	(&Response{
		Status:     true,
		Message:    "trace results",
		StatusCode: http.StatusOK,
		Data:       results,
	}).Success(w)
}

func parseSizeFromRequest(r *http.Request) (int, int) {
	var err error

	var width int

	var height int

	if r.URL.Query().Get("width") == "" || r.URL.Query().Get("height") == "" {
		width = defaultScreenWidth
		height = defaultScreenHeight
	} else if width, err = strconv.Atoi(r.URL.Query().Get("width")); err != nil {
		width = defaultScreenWidth
		height = defaultScreenHeight
	} else if height, err = strconv.Atoi(r.URL.Query().Get("height")); err != nil {
		width = defaultScreenWidth
		height = defaultScreenHeight
	}

	return width, height
}
