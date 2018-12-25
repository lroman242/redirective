package tracer

import (
	"net/http"
	"net/url"
)

type redirect struct {
	From            *url.URL               `json:"from"`
	To              *url.URL               `json:"to"`
	RequestHeaders  *http.Header           `json:"request_headers"`
	ResponseHeaders *http.Header           `json:"response_headers"`
	Cookies         []*http.Cookie         `json:"cookies"`
	Body            []byte                 `json:"body"`
	Status          int                    `json:"status"`
	Initiator       string                 `json:"initiator"`
	OtherInfo       map[string]interface{} `json:"other_info"`
}

func NewRedirect(from, to *url.URL, requestHeaders, responseHeaders *http.Header, cookies []*http.Cookie, body []byte, status int, initiator string) *redirect {
	return &redirect{
		From:            from,
		To:              to,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
		Cookies:         cookies,
		Body:            body,
		Status:          status,
		Initiator:       initiator,
	}
}
