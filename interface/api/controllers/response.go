package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

const notConvertedToJSONErrorMessage = "response data cannot be converted to json"

// Response type describe common http response
type Response struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

// Success should be used to send success response (200 http status code)
func (r *Response) Success(w http.ResponseWriter) {
	r.Status = true

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}

	jsonResponse, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error while json encode: %q", err.Error())

		r.StatusCode = http.StatusInternalServerError
		r.Message = notConvertedToJSONErrorMessage
		r.Data = nil
		r.Failed(w)

		return
	}

	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResponse)
}

// Failed should be used to send error response (40x or 50x http status code)
func (r *Response) Failed(w http.ResponseWriter) {
	r.Status = false

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusBadRequest
	}

	jsonResponse, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error while json encode: %q", err.Error())

		r.StatusCode = http.StatusInternalServerError
		r.Message = notConvertedToJSONErrorMessage
		r.Data = nil
		r.Failed(w)

		return
	}

	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResponse)
}
