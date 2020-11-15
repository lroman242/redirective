package controllers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestTraceController_FindTraceResults(t *testing.T) {
	testResultsID := "SomeTestResultsID"
	expectedResponse := `{"status":true,"message":"trace results","status_code":200,"data":{"id":"SomeTestResultsID","redirects":null,"screenshot":"some/path/to/screenshot.png","url":"http://example.domain/screenshots/screenshot.png"}}`

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceResults := &domain.TraceResults{
		ID:         testResultsID,
		Redirects:  nil,
		Screenshot: "some/path/to/screenshot.png",
		URL:        "http://example.domain/screenshots/screenshot.png",
	}

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().FindTrace(testResultsID).Times(1).Return(traceResults, nil)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/find/"+testResultsID, nil)
	params := httprouter.Params{struct {
		Key   string
		Value string
	}{Key: "id", Value: testResultsID}}

	controller := NewTraceController(traceInteractor, "screenshots/", logger)
	controller.FindTraceResults(response, request, params)

	if response.Code != http.StatusOK {
		t.Error("wrong http status code received. expected 200")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_FindTraceResults_TraceInteractor_FindTrace_Error(t *testing.T) {
	testResultsID := "SomeTestResultsID"
	expectedResponse := `{"status":false,"message":"an error occurred. error: some error","status_code":500,"data":null}`
	expectedError := errors.New("some error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().FindTrace(testResultsID).Times(1).Return(nil, expectedError)
	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError).Times(1)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/find/"+testResultsID, nil)
	params := httprouter.Params{struct {
		Key   string
		Value string
	}{Key: "id", Value: testResultsID}}

	controller := NewTraceController(traceInteractor, "screenshots/", logger)
	controller.FindTraceResults(response, request, params)

	if response.Code != http.StatusInternalServerError {
		t.Error("wrong http status code received. expected 500")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_TraceUrl(t *testing.T) {
	expectedResponse := `{"status":true,"message":"url traced","status_code":200,"data":{"id":"SomeTestResultsID","redirects":null,"screenshot":"some/path/to/screenshot.png","url":"http://example.domain/screenshots/screenshot.png"}}`

	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)

	expectedScreenshotsPath := "screenshots/"
	testResultsID := "SomeTestResultsID"

	expectedTraceResults := &domain.TraceResults{
		ID:         testResultsID,
		Redirects:  nil,
		Screenshot: "some/path/to/screenshot.png",
		URL:        "http://example.domain/screenshots/screenshot.png",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Trace(expectedURL, expectedScreenshotsPath).Times(1).Return(expectedTraceResults, nil)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/trace?url=" + expectedStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.TraceUrl(response, request, params)

	if response.Code != http.StatusOK {
		t.Error("wrong http status code received. expected 200")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_TraceUrl_InvalidURL_Error(t *testing.T) {
	inputStrURL := "invalid_url_input"
	expectedResponse := `{"status":false,"message":"invalid url parse \"invalid_url_input\": invalid URI for request","status_code":400,"data":null}`
	expectedScreenshotsPath := "screenshots/"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/trace?url=" + inputStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.TraceUrl(response, request, params)

	if response.Code != http.StatusBadRequest {
		t.Error("wrong http status code received. expected 400")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_TraceUrl_NoUrlProvided_Error(t *testing.T) {
	inputStrURL := ""
	expectedResponse := `{"status":false,"message":"url parameter is required","status_code":400,"data":null}`
	expectedScreenshotsPath := "screenshots/"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/trace?url=" + inputStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.TraceUrl(response, request, params)

	if response.Code != http.StatusBadRequest {
		t.Error("wrong http status code received. expected 400")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_TraceUrl_TraceInteractor_Trace_Error(t *testing.T) {
	expectedResponse := `{"status":false,"message":"an error occurred. error: some error","status_code":500,"data":null}`

	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)

	expectedScreenshotsPath := "screenshots/"

	expectedError := errors.New("some error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Trace(expectedURL, expectedScreenshotsPath).Times(1).Return(nil, expectedError)
	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError).Times(1)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/trace?url=" + expectedStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.TraceUrl(response, request, params)

	if response.Code != http.StatusInternalServerError {
		t.Error("wrong http status code received. expected 500")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_Screenshot(t *testing.T) {
	expectedResponse := `{"status":true,"message":"url traced","status_code":200,"data":"some/path/to/screenshot.png"}`

	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)
	expectedWidth := 1920
	expectedHeight := 1080
	expectedScreenshotsPath := "screenshots/"

	expectedScreenshotResults := "some/path/to/screenshot.png"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Screenshot(expectedURL, expectedWidth, expectedHeight, expectedScreenshotsPath).Times(1).Return(expectedScreenshotResults, nil)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + expectedStrURL + "&width=" + strconv.Itoa(expectedWidth) + "&height=" + strconv.Itoa(expectedHeight), nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)

	if response.Code != http.StatusOK {
		t.Error("wrong http status code received. expected 200")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_Screenshot_NoUrlProvided_Error(t *testing.T) {
	inputStrURL := ""
	expectedResponse := `{"status":false,"message":"url parameter is required","status_code":400,"data":null}`
	expectedScreenshotsPath := "screenshots/"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + inputStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)

	if response.Code != http.StatusBadRequest {
		t.Error("wrong http status code received. expected 400")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_Screenshot_InvalidURL_Error(t *testing.T) {
	inputStrURL := "invalid_url_input"
	expectedResponse := `{"status":false,"message":"invalid url parse \"invalid_url_input\": invalid URI for request","status_code":400,"data":null}`
	expectedScreenshotsPath := "screenshots/"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + inputStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)

	if response.Code != http.StatusBadRequest {
		t.Error("wrong http status code received. expected 400")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_Screenshot_TraceInteractor_Screenshot_Error(t *testing.T) {
	expectedResponse := `{"status":false,"message":"an error occurred. error: some error","status_code":500,"data":null}`

	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)
	expectedWidth := 1920
	expectedHeight := 1080
	expectedScreenshotsPath := "screenshots/"

	expectedError := errors.New("some error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Screenshot(expectedURL, expectedWidth, expectedHeight, expectedScreenshotsPath).Times(1).Return("", expectedError)
	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError).Times(1)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + expectedStrURL + "&width=" + strconv.Itoa(expectedWidth) + "&height=" + strconv.Itoa(expectedHeight), nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)

	if response.Code != http.StatusInternalServerError {
		t.Error("wrong http status code received. expected 500")
	}
	if response.Body.String() != expectedResponse {
		t.Error("wrong response body received")
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong response headers. expected `content-type: application/json` header")
	}
}

func TestTraceController_Screenshot_Width_And_Height_NotProvided(t *testing.T) {
	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)
	expectedWidth := 1920
	expectedHeight := 1080
	expectedScreenshotsPath := "screenshots/"
	expectedScreenshotResults := "some/path/to/screenshot.png"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Screenshot(expectedURL, expectedWidth, expectedHeight, expectedScreenshotsPath).Times(1).Return(expectedScreenshotResults, nil)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + expectedStrURL, nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)
}

func TestTraceController_Screenshot_Width_NotProvided(t *testing.T) {
	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)
	expectedWidth := 1920
	expectedHeight := 1080
	expectedScreenshotsPath := "screenshots/"
	expectedScreenshotResults := "some/path/to/screenshot.png"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Screenshot(expectedURL, expectedWidth, expectedHeight, expectedScreenshotsPath).Times(1).Return(expectedScreenshotResults, nil)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + expectedStrURL + "&width=&height=300", nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)
}

func TestTraceController_Screenshot_Height_NotProvided(t *testing.T) {
	expectedStrURL := "http://google.com.ua"
	expectedURL, _ := url.ParseRequestURI(expectedStrURL)
	expectedWidth := 1920
	expectedHeight := 1080
	expectedScreenshotsPath := "screenshots/"
	expectedScreenshotResults := "some/path/to/screenshot.png"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	traceInteractor := mocks.NewMockTraceInteractor(mockCtrl)
	traceInteractor.EXPECT().Screenshot(expectedURL, expectedWidth, expectedHeight, expectedScreenshotsPath).Times(1).Return(expectedScreenshotResults, nil)
	logger := mocks.NewMockLogger(mockCtrl)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/screenshot?url=" + expectedStrURL + "&width=500&height=", nil)
	params := make(httprouter.Params, 0,0)

	controller := NewTraceController(traceInteractor, expectedScreenshotsPath, logger)
	controller.Screenshot(response, request, params)
}