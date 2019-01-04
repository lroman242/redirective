package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func (r *Response) Success(w http.ResponseWriter) {
	r.Status = true

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}

	jsonResponse, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error while json encode: %q", err.Error())
		r.StatusCode = http.StatusInternalServerError
		r.Message = "response data cannot be converted to json"
		r.Data = nil
		r.Failed(w)
		return
	}

	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (r *Response) Failed(w http.ResponseWriter) {
	r.Status = false

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusBadRequest
	}

	jsonResponse, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error while json encode: %q", err.Error())
		r.StatusCode = http.StatusInternalServerError
		r.Message = "response data cannot be converted to json"
		r.Data = nil
		r.Failed(w)
		return
	}

	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
