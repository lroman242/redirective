// Package controllers implements methods to handle http requests
package controllers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ChromeScreenshot function create image (screenshot) of active browser tab
func ChromeScreenshot(w http.ResponseWriter, r *http.Request, screenshotsStoragePath string) {
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

	chr := tracer.NewChromeTracer(remote, parseScreenSizeFromRequest(r), screenshotsStoragePath)

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

	screenShotFileName := randomScreenshotFileName()

	err = chr.Screenshot(targetURL, parseScreenSizeFromRequest(r), screenShotFileName)
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
		Data:       screenShotFileName}).Success(w)
}

// ChromeTrace parse a trace path for provided url
func ChromeTrace(w http.ResponseWriter, r *http.Request, screenshotsStoragePath string, col *mongo.Collection) {
	screenShotFileName := randomScreenshotFileName()
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
		if err = remote.Close(); err != nil {
			log.Printf("remote.Close error: %s \n", err)
		}
	}()
	// create new tracer instance
	chr := tracer.NewChromeTracer(remote, parseScreenSizeFromRequest(r), screenshotsStoragePath)
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
	redirects, err := chr.Trace(targetURL, screenShotFileName)
	if err != nil {
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprintf("sorry, an error occurred. %s", err),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}

	jsonRedirects := tracer.NewJSONRedirects(redirects)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := col.InsertOne(ctx, bson.M{"redirects": jsonRedirects, "screenshot": screenShotFileName})
	if err != nil {
		log.Printf("error occurred during saving trace results. error: %s \n", err)

		(&response.Response{
			Status:     true,
			Message:    "url successfully traced",
			StatusCode: 200,
			Data: struct {
				Redirects  []*tracer.JSONRedirect `json:"redirects"`
				Screenshot string                 `json:"screenshot"`
			}{
				Redirects:  jsonRedirects,
				Screenshot: screenShotFileName,
			}}).Success(w)

		return
	}

	(&response.Response{
		Status:     true,
		Message:    "url successfully traced",
		StatusCode: 200,
		Data: struct {
			Redirects  []*tracer.JSONRedirect `json:"redirects"`
			Screenshot string                 `json:"screenshot"`
			ID         interface{}            `json:"id"`
		}{
			Redirects:  jsonRedirects,
			Screenshot: screenShotFileName,
			ID:         res.InsertedID,
		}}).Success(w)
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

func LoadTraceResults(w http.ResponseWriter, r *http.Request, col *mongo.Collection, id string) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid ObjectID. error: %s \n", err)
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprint("sorry, an error occurred. invalid id"),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}

	trace := new(struct {
		Redirects  []*tracer.JSONRedirect `json:"redirects"`
		Screenshot string                 `json:"screenshot"`
	})

	err = col.FindOne(ctx, bson.M{"_id": ID}).Decode(trace)
	if err != nil {
		log.Printf("mongodb reulsts decode failed. error: %s \n", err)
		(&response.Response{
			Status:     false,
			Message:    fmt.Sprint("sorry, an error occurred. trace not found"),
			StatusCode: 500,
			Data:       nil}).Failed(w)

		return
	}

	log.Printf("%+v\n", trace)

	(&response.Response{
		Status:     true,
		Message:    "url successfully traced",
		StatusCode: 200,
		Data:       trace}).Success(w)
}
