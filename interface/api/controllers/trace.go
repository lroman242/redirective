package controllers

import (
	"fmt"
	"github.com/lroman242/redirective/usecase/interactor"
	"net/http"
	"net/url"
)

type TraceController interface {
	TraceUrl(http.ResponseWriter, *http.Request)
}

type traceController struct {
	traceInteractor interactor.TraceInteractor
	assetsPath      string
}

func NewTraceController(ti interactor.TraceInteractor, assetsPath string) TraceController {
	return &traceController{ti, assetsPath}
}

func (tc *traceController) TraceUrl(w http.ResponseWriter, r *http.Request) {
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: 400,
			Data:       nil}).Failed(w)

		return
	}

	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: 400,
			Data:       nil}).Failed(w)

		return
	}

	results, err := tc.traceInteractor.Trace(targetURL, tc.assetsPath)

	(&Response{
		Status:     true,
		Message:    "url traced",
		StatusCode: 200,
		Data:       results,
	}).Success(w)
}
