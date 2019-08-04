package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse_Failed(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	response := &Response{}

	response.Failed(responseWriter)

	if responseWriter.Code != http.StatusBadRequest {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusNotFound, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong content-type header value")
	}

	body := responseWriter.Body

	unmarshalResponse := &Response{}

	readBuf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(readBuf, &unmarshalResponse)
	if err != nil {
		t.Error(err)
	}

	if unmarshalResponse.Status != false {
		t.Error("expect Response.Status equal to false")
	}

	if unmarshalResponse.Message != "" {
		t.Error("expect Response.Message equal to empty string")
	}

	if unmarshalResponse.StatusCode != http.StatusBadRequest {
		t.Error("expect Response.StatusCode equal to 400")
	}

	if unmarshalResponse.Data != nil {
		t.Error("expect Response.Data equal to nil")
	}
}

func TestResponse_Failed_With_Data(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	response := &Response{true, "Message", http.StatusNotFound, []string{"test", "strings"}}

	response.Failed(responseWriter)

	if responseWriter.Code != http.StatusNotFound {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusNotFound, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong content-type header value")
	}

	body := responseWriter.Body

	unmarshalResponse := &Response{}

	readBuf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(readBuf, &unmarshalResponse)
	if err != nil {
		t.Error(err)
	}

	if unmarshalResponse.Status != false {
		t.Error("expect Response.Status equal to false")
	}

	if unmarshalResponse.Message != "Message" {
		t.Error("expect Response.Message equal to `Message`")
	}

	if unmarshalResponse.StatusCode != http.StatusNotFound {
		t.Error("expect Response.StatusCode equal to 404")
	}

	if unmarshalResponse.Data == nil {
		t.Error("expect Response.Data not equal to nil")
	}

	interfaceData := unmarshalResponse.Data.([]interface{})
	if len(interfaceData) != 2 {
		t.Error("expect Response.Data contains 2 elements")
	}
	if interfaceData[0].(string) != "test" {
		t.Error("expect Response.Data[0] equal to `test`")
	}
	if interfaceData[1].(string) != "strings" {
		t.Error("expect Response.Data[1] equal to `strings`")
	}
}

func TestResponse_Success(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	response := &Response{}

	response.Success(responseWriter)

	if responseWriter.Code != http.StatusOK {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusOK, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong content-type header value")
	}

	body := responseWriter.Body

	unmarshalResponse := &Response{}

	readBuf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(readBuf, &unmarshalResponse)
	if err != nil {
		t.Error(err)
	}

	if unmarshalResponse.Status != true {
		t.Error("expect Response.Status equal to true")
	}

	if unmarshalResponse.Message != "" {
		t.Error("expect Response.Message equal to empty string")
	}

	if unmarshalResponse.StatusCode != http.StatusOK {
		t.Error("expect Response.StatusCode equal to 200")
	}

	if unmarshalResponse.Data != nil {
		t.Error("expect Response.Data equal to nil")
	}
}

func TestResponse_Success_With_Data(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	response := &Response{false, "SuccessMessage", http.StatusAccepted, []string{"success", "strings"}}

	response.Success(responseWriter)

	if responseWriter.Code != http.StatusAccepted {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusAccepted, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong content-type header value")
	}

	body := responseWriter.Body

	unmarshalResponse := &Response{}

	readBuf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(readBuf, &unmarshalResponse)
	if err != nil {
		t.Error(err)
	}

	if unmarshalResponse.Status != true {
		t.Error("expect Response.Status equal to true")
	}

	if unmarshalResponse.Message != "SuccessMessage" {
		t.Error("expect Response.Message equal to `SuccessMessage`")
	}

	if unmarshalResponse.StatusCode != http.StatusAccepted {
		t.Error("expect Response.StatusCode equal to 202")
	}

	if unmarshalResponse.Data == nil {
		t.Error("expect Response.Data not equal to nil")
	}

	interfaceData := unmarshalResponse.Data.([]interface{})
	if len(interfaceData) != 2 {
		t.Error("expect Response.Data contains 2 elements")
	}
	if interfaceData[0].(string) != "success" {
		t.Error("expect Response.Data[0] equal to `success`")
	}
	if interfaceData[1].(string) != "strings" {
		t.Error("expect Response.Data[1] equal to `strings`")
	}
}

func TestResponse_Success_With_Invalid_Data(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	data := make(chan int)
	response := &Response{false, "SuccessMessage", http.StatusAccepted, data}

	response.Success(responseWriter)

	if responseWriter.Code != http.StatusInternalServerError {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusInternalServerError, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong content-type header value")
	}

	body := responseWriter.Body

	unmarshalResponse := &Response{}

	readBuf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(readBuf, &unmarshalResponse)
	if err != nil {
		t.Error(err)
	}

	if unmarshalResponse.Status != false {
		t.Error("expect Response.Status equal to false")
	}

	if unmarshalResponse.Message != "response data cannot be converted to json" {
		t.Error("expect Response.Message equal to `response data cannot be converted to json`")
	}

	if unmarshalResponse.StatusCode != http.StatusInternalServerError {
		t.Error("expect Response.StatusCode equal to 500")
	}

	if unmarshalResponse.Data != nil {
		t.Error("expect Response.Data equal to nil")
	}
}

func TestResponse_Failed_With_Invalid_Data(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	data := make(chan int)
	response := &Response{true, "FailedMessage", http.StatusAccepted, data}

	response.Failed(responseWriter)

	if responseWriter.Code != http.StatusInternalServerError {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusInternalServerError, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Error("wrong content-type header value")
	}

	body := responseWriter.Body

	unmarshalResponse := &Response{}

	readBuf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(readBuf, &unmarshalResponse)
	if err != nil {
		t.Error(err)
	}

	if unmarshalResponse.Status != false {
		t.Error("expect Response.Status equal to false")
	}

	if unmarshalResponse.Message != "response data cannot be converted to json" {
		t.Error("expect Response.Message equal to `response data cannot be converted to json`")
	}

	if unmarshalResponse.StatusCode != http.StatusInternalServerError {
		t.Error("expect Response.StatusCode equal to 500")
	}

	if unmarshalResponse.Data != nil {
		t.Error("expect Response.Data equal to nil")
	}
}
