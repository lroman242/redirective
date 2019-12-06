// Package controllers implements methods to handle http requests
package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/lroman242/redirective/response"
	"github.com/lroman242/redirective/tracer"
	"github.com/raff/godet"
)

const defaultScreenWidth = 1920
const defaultScreenHeight = 1080

const screenshotsSavePath = "assets/screenshots/"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ChromeScreenshot function create image (screenshot) of active browser tab
func ChromeScreenshot(w http.ResponseWriter, r *http.Request) {
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("cannot connect to Chrome instance: %s", err),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}

	defer func() {
		err = remote.Close()
		if err != nil {
			log.Printf("remote.Close error: %s", err)
		}
	}()

	chr := tracer.NewChromeTracer(remote)

	urlToTrace := r.URL.Query().Get("url")
	if urlToTrace == "" {
		(&response.Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: 400,
			Data:       nil}).Failed(w)

		return
	}
	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: 400,
			Data:       nil}).Failed(w)

		return
	}

	screenShotPath := screenshotsSavePath + randomScreenshotFileName()

	err = chr.Screenshot(targetURL, parseScreenSizeFromRequest(r), screenShotPath)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("an error occurred. %s", err),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}

	(&response.Response{
		Status:     true,
		Message:    "url successfully traced",
		StatusCode: 200,
		Data:       screenShotPath}).Success(w)
}

// ChromeTrace parse a trace path for provided url
func ChromeTrace(w http.ResponseWriter, r *http.Request) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("cannot connect to Chrome instance: %s", err),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}
	// close connection
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
		(&response.Response{
			Status:     false,
			Message:    "url parameter is required",
			StatusCode: 400,
			Data:       nil}).Failed(w)

		return
	}
	// convert raw url string to url.URL
	targetURL, err := url.ParseRequestURI(urlToTrace)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("invalid url %s", err),
			StatusCode: 400,
			Data:       nil}).Failed(w)

		return
	}
	// process tracing
	redirects, err := chr.Trace(targetURL)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("sorry, an error occurred. %s", err),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}

	(&response.Response{
		Status:     true,
		Message:    "url successfully traced",
		StatusCode: 200,
		Data:       tracer.NewJSONRedirects(redirects)}).Success(w)
}

// parseScreenSizeFromRequest - parse screen width and height from request or use default values
func parseScreenSizeFromRequest(r *http.Request) *tracer.ScreenSize {
	var width int

	widthStr := r.URL.Query().Get("width")
	if widthStr == "" {
		widthStr = strconv.Itoa(defaultScreenWidth)
	}

	width, err := strconv.Atoi(widthStr)
	if err != nil {
		width = defaultScreenWidth
	}

	var height int

	heightStr := r.URL.Query().Get("height")
	if widthStr == "" {
		heightStr = strconv.Itoa(defaultScreenHeight)
	}

	height, err = strconv.Atoi(heightStr)
	if err != nil {
		height = defaultScreenHeight
	}

	return tracer.NewScreenSize(width, height)
}

func randomScreenshotFileName() string {
	b := make([]byte, 16)

	for i := range b {
		b[i] = charset[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}

	return string(b) + `.png`
}
