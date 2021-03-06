package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const jsonMimeType = "application/json"
const stringTestValue1 = "testVal1"
const stringTestValue2 = "testVal2"
const stringTestMessage = "Message"

func TestResponse_Failed(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	response := &Response{}

	response.Failed(responseWriter)

	if responseWriter.Code != http.StatusBadRequest {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusNotFound, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != jsonMimeType {
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
	response := &Response{true, stringTestMessage, http.StatusNotFound, []string{stringTestValue1, stringTestValue2}}

	response.Failed(responseWriter)

	if responseWriter.Code != http.StatusNotFound {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusNotFound, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != jsonMimeType {
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
	} else if unmarshalResponse.Message != stringTestMessage {
		t.Errorf("expect Response.Message equal to `%s`", stringTestMessage)
	} else if unmarshalResponse.StatusCode != http.StatusNotFound {
		t.Error("expect Response.StatusCode equal to 404")
	} else if unmarshalResponse.Data == nil {
		t.Error("expect Response.Data not equal to nil")
	}

	interfaceData := unmarshalResponse.Data.([]interface{})
	if len(interfaceData) != 2 {
		t.Error("expect Response.Data contains 2 elements")
	} else if interfaceData[0].(string) != stringTestValue1 {
		t.Errorf("expect Response.Data[0] equal to `%s`", stringTestValue1)
	} else if interfaceData[1].(string) != stringTestValue2 {
		t.Errorf("expect Response.Data[1] equal to `%s`", stringTestValue2)
	}
}

func TestResponse_Success(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	response := &Response{}

	response.Success(responseWriter)

	if responseWriter.Code != http.StatusOK {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusOK, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != jsonMimeType {
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
	response := &Response{false, stringTestMessage, http.StatusAccepted, []string{stringTestValue1, stringTestValue2}}

	response.Success(responseWriter)

	if responseWriter.Code != http.StatusAccepted {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusAccepted, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != jsonMimeType {
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

	if unmarshalResponse.Message != stringTestMessage {
		t.Errorf("expect Response.Message equal to `%s`", stringTestMessage)
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
	} else if interfaceData[0].(string) != stringTestValue1 {
		t.Errorf("expect Response.Data[0] equal to `%s`", stringTestValue1)
	} else if interfaceData[1].(string) != stringTestValue2 {
		t.Errorf("expect Response.Data[1] equal to `%s`", stringTestValue2)
	}
}

func TestResponse_Success_With_Invalid_Data(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	data := make(chan int)
	response := &Response{false, stringTestMessage, http.StatusAccepted, data}

	response.Success(responseWriter)

	if responseWriter.Code != http.StatusInternalServerError {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusInternalServerError, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != jsonMimeType {
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

	if unmarshalResponse.Message != notConvertedToJSONErrorMessage {
		t.Errorf("expect Response.Message equal to `%s`", notConvertedToJSONErrorMessage)
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
	response := &Response{true, stringTestMessage, http.StatusAccepted, data}

	response.Failed(responseWriter)

	if responseWriter.Code != http.StatusInternalServerError {
		t.Errorf("wrong response status code. expect %d but get %d", http.StatusInternalServerError, responseWriter.Code)
	}

	if responseWriter.Header().Get("Content-Type") != jsonMimeType {
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

	if unmarshalResponse.Message != notConvertedToJSONErrorMessage {
		t.Errorf("expect Response.Message equal to `%s`", notConvertedToJSONErrorMessage)
	}

	if unmarshalResponse.StatusCode != http.StatusInternalServerError {
		t.Error("expect Response.StatusCode equal to 500")
	}

	if unmarshalResponse.Data != nil {
		t.Error("expect Response.Data equal to nil")
	}
}
