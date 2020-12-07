package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/lroman242/redirective/usecase/interactor"
	"net/http"
	"net/url"
	"strconv"
)

const defaultScreenWidth = 1920
const defaultScreenHeight = 1080

type TraceController interface {
	TraceURL(http.ResponseWriter, *http.Request, httprouter.Params)
	Screenshot(http.ResponseWriter, *http.Request, httprouter.Params)
	FindTraceResults(http.ResponseWriter, *http.Request, httprouter.Params)
}

type traceController struct {
	traceInteractor interactor.TraceInteractor
	assetsPath      string
	log             logger.Logger
}

// NewTraceController
func NewTraceController(ti interactor.TraceInteractor, assetsPath string, log logger.Logger) TraceController {
	return &traceController{ti, assetsPath, log}
}

func (tc *traceController) TraceURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: http.StatusBadRequest,
			Data:       nil}).Failed(w)

		return
	}

	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: http.StatusBadRequest,
			Data:       nil}).Failed(w)

		return
	}

	results, err := tc.traceInteractor.Trace(targetURL, tc.assetsPath)
	if err != nil {
		tc.log.Error(err)
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("an error occurred. error: %s", err),
			StatusCode: http.StatusInternalServerError,
			Data:       nil}).Failed(w)

		return
	}

	(&Response{
		Status:     true,
		Message:    "url traced",
		StatusCode: http.StatusOK,
		Data:       results,
	}).Success(w)
}

func (tc *traceController) Screenshot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: http.StatusBadRequest,
			Data:       nil}).Failed(w)

		return
	}

	var err error

	var width int

	var height int

	if r.URL.Query().Get("width") == "" || r.URL.Query().Get("height") == "" {
		width = defaultScreenWidth
		height = defaultScreenHeight
	} else {
		width, err = strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil {
			width = defaultScreenWidth
			height = defaultScreenHeight
		} else {
			height, err = strconv.Atoi(r.URL.Query().Get("height"))
			if err != nil {
				width = defaultScreenWidth
				height = defaultScreenHeight
			}
		}
	}

	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: http.StatusBadRequest,
			Data:       nil}).Failed(w)

		return
	}

	screenshotURL, err := tc.traceInteractor.Screenshot(targetURL, width, height, tc.assetsPath)
	if err != nil {
		tc.log.Error(err)
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("an error occurred. error: %s", err),
			StatusCode: http.StatusInternalServerError,
			Data:       nil}).Failed(w)

		return
	}

	(&Response{
		Status:     true,
		Message:    "url traced",
		StatusCode: http.StatusOK,
		Data:       screenshotURL,
	}).Success(w)
}

func (tc *traceController) FindTraceResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	results, err := tc.traceInteractor.FindTrace(id)
	if err != nil {
		tc.log.Error(err)
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("an error occurred. error: %s", err),
			StatusCode: http.StatusInternalServerError,
			Data:       nil}).Failed(w)

		return
	}

	(&Response{
		Status:     true,
		Message:    "trace results",
		StatusCode: http.StatusOK,
		Data:       results,
	}).Success(w)
}
