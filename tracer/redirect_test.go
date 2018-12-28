package tracer

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestNewRedirect(t *testing.T) {
	from, _ := url.Parse("http://google.com")
	to, _ := url.Parse("https://google.com")

	requestHeaders := &http.Header{}
	requestHeaders.Set("requestTestKey", "requestTestVal1")

	responseHeaders := &http.Header{}
	responseHeaders.Set("responseTestKey", "responseTestVal1")

	cookies := []*http.Cookie{&http.Cookie{
		Name:     "testCookie",
		Domain:   "google.com",
		Expires:  time.Now(),
		HttpOnly: false,
		MaxAge:   5,
		Path:     "/",
		Secure:   false,
		Value:    "testCookieValue",
	}}

	//body := []byte("test body")

	status := http.StatusFound

	initiator := "test script"

	redirect := NewRedirect(from, to, requestHeaders, responseHeaders, cookies, /*body,*/ status, initiator)

	if redirect.Status != status {
		t.Error("Invalid status on redirect creating")
	}

	if redirect.From.String() != "http://google.com" {
		t.Error("Invalid From URL on redirect creating")
	}

	if redirect.To.String() != "https://google.com" {
		t.Error("Invalid To URL on redirect creating")
	}

	//if string(redirect.Body) != string(body) {
	//	t.Error("Invalid Body on redirect creating")
	//}

	if redirect.Initiator != initiator {
		t.Error("Invalid Initiator on redirect creating")
	}

	if len(redirect.Cookies) <= 0 {
		t.Error("No Cookies on redirect creating")
	}

	if cookie := redirect.Cookies[0]; cookie.Name != "testCookie" {
		t.Error("Wrong Cookie name on redirect creating")
	}
	if cookie := redirect.Cookies[0]; cookie.Value != "testCookieValue" {
		t.Error("Wrong Cookie value on redirect creating")
	}

	if headerVal := redirect.RequestHeaders.Get("requestTestKey"); headerVal != "requestTestVal1" {
		t.Errorf("Wrong Request header value. Expect %s but get %s", "requestTestVal1", headerVal)
	}
	if headerVal := redirect.ResponseHeaders.Get("responseTestKey"); headerVal != "responseTestVal1" {
		t.Errorf("Wrong Response header value. Expect %s but get %s", "responseTestVal1", headerVal)
	}
}
