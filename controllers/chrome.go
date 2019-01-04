package controllers

import (
	"fmt"
	"github.com/lroman242/redirective/response"
	"github.com/lroman242/redirective/tracer"
	"github.com/raff/godet"
	"log"
	"net/http"
	"net/url"
)

func ChromeScreenshot(w http.ResponseWriter, r *http.Request)  {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		(&response.Response{false, fmt.Sprintf("cannot connect to Chrome instance: %s", err), 500, nil}).Failed(w)
		return
	}

	defer func() {
		err = remote.Close()
		if err != nil {
			log.Printf("remote.Close error: %s", err)
		}
	}()

	// create new tracer instance
	chr := tracer.NewChromeTracer(remote)

	// check url
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&response.Response{false, "url parameter is required", 400, nil}).Failed(w)
		return
	}

	// convert raw url string to url.URL
	targetUrl, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&response.Response{false, fmt.Sprintf("invalid url %s", err), 400, nil}).Failed(w)
		return
	}

	//TODO: unique image name
	//TODO: custom screen size
	screenShotPath := "assets/screenshots/test.png"
	err = chr.Screenshot(targetUrl, tracer.NewScreenSize(1920,1080), screenShotPath)
	if err != nil {
		(&response.Response{false, fmt.Sprintf("an error occurred. %s", err), 500, nil}).Failed(w)
		return
	}

	(&response.Response{true, "url successfully traced", 200, screenShotPath}).Success(w)
}

func ChromeTrace(w http.ResponseWriter, r *http.Request) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		(&response.Response{false, fmt.Sprintf("cannot connect to Chrome instance: %s", err), 500, nil}).Failed(w)
		return
	}

	defer func() {
		err = remote.Close()
		if err != nil {
			log.Printf("remote.Close error: %s", err)
		}
	}()

	// create new tracer instance
	chr := tracer.NewChromeTracer(remote)

	// check url
	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&response.Response{false, "url parameter is required", 400, nil}).Failed(w)
		return
	}

	// convert raw url string to url.URL
	targetUrl, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&response.Response{false, fmt.Sprintf("invalid url %s", err), 400, nil}).Failed(w)
		return
	}

	// process tracing
	redirects, err := chr.GetTrace(targetUrl)
	if err != nil {
		(&response.Response{false, fmt.Sprintf("sorry, an error occurred. %s", err), 500, nil}).Failed(w)
		return
	}

	(&response.Response{true, "url successfully traced", 200, tracer.NewJSONRedirects(redirects)}).Success(w)
}
