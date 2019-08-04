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

	redirect := NewRedirect(from, to, requestHeaders, responseHeaders, cookies /*body,*/, status, initiator)

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

func TestNewJSONCookie(t *testing.T) {
	header := http.Header{}
	header.Add("Set-Cookie", "foo1=bar1; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test1.com; Path=/;")
	header.Add("Set-Cookie", "foo2=bar2; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test2.com; Path=/;")
	cookies := (&http.Response{Header: header}).Cookies()

	jsonCookies := NewJSONCookies(cookies)

	if len(jsonCookies) != 2 {
		t.Errorf("expect to get 2 cookies but get %d", len(jsonCookies))
	}

	if jsonCookies[0].Raw != cookies[0].Raw {
		t.Error("Invalid 1st Cookie raw value")
	}
	if jsonCookies[1].Raw != cookies[1].Raw {
		t.Error("Invalid 2nd Cookie raw value")
	}
}

func TestNewJSONRedirects(t *testing.T) {
	header1 := http.Header{}
	header1.Add("Set-Cookie", "foo1=bar1; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test1.com; Path=/;")
	header1.Add("Set-Cookie", "foo2=bar2; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test2.com; Path=/;")
	cookies1 := (&http.Response{Header: header1}).Cookies()

	header2 := http.Header{}
	header2.Add("Set-Cookie", "foo3=bar3; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test3.com; Path=/;")
	header2.Add("Set-Cookie", "foo4=bar4; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test4.com; Path=/;")
	cookies2 := (&http.Response{Header: header1}).Cookies()

	var redirects []*redirect

	fromUrl1, _ := url.ParseRequestURI("http://google.com")
	toUrl1, _ := url.ParseRequestURI("https://google.com")
	fromUrl2, _ := url.ParseRequestURI("https://google.com")
	toUrl2, _ := url.ParseRequestURI("https://www.google.com")
	redirects = append(redirects, NewRedirect(fromUrl1, toUrl1, &header1, &header2, cookies1, 303, "other"))
	redirects = append(redirects, NewRedirect(fromUrl2, toUrl2, &header2, &header1, cookies2, 303, "other"))

	jsonRedirects := NewJSONRedirects(redirects)

	if len(jsonRedirects) != len(redirects) {
		t.Errorf("wrong redirects amount. expect %d but get %d", len(redirects), len(jsonRedirects))
	}

	if jsonRedirects[0].Status != redirects[0].Status {
		t.Errorf("wrong Status value. expect %d but get %d", redirects[0].Status, jsonRedirects[0].Status)
	}

	if jsonRedirects[0].From != redirects[0].From.String() {
		t.Errorf("wrong From value. expect %s but get %s", redirects[0].From.String(), jsonRedirects[0].From)
	}

	if jsonRedirects[0].To != redirects[0].To.String() {
		t.Errorf("wrong To value. expect %s but get %s", redirects[0].To.String(), jsonRedirects[0].To)
	}

	if jsonRedirects[0].Initiator != redirects[0].Initiator {
		t.Errorf("wrong Initiator value. expect %s but get %s", redirects[0].Initiator, jsonRedirects[0].Initiator)
	}

	if jsonRedirects[0].Cookies[0].Raw != redirects[0].Cookies[0].Raw {
		t.Errorf("wrong Cookie[0] value. expect %s but get %s", redirects[0].Cookies[0].Raw, jsonRedirects[0].Cookies[0].Raw)
	}

	if jsonRedirects[0].Cookies[1].Raw != redirects[0].Cookies[1].Raw {
		t.Errorf("wrong Cookie[1] value. expect %s but get %s", redirects[0].Cookies[1].Raw, jsonRedirects[0].Cookies[1].Raw)
	}

	if jsonRedirects[1].Status != redirects[1].Status {
		t.Errorf("wrong Status value. expect %d but get %d", redirects[1].Status, jsonRedirects[1].Status)
	}

	if jsonRedirects[1].From != redirects[1].From.String() {
		t.Errorf("wrong From value. expect %s but get %s", redirects[1].From.String(), jsonRedirects[1].From)
	}

	if jsonRedirects[1].To != redirects[1].To.String() {
		t.Errorf("wrong To value. expect %s but get %s", redirects[1].To.String(), jsonRedirects[1].To)
	}

	if jsonRedirects[1].Initiator != redirects[1].Initiator {
		t.Errorf("wrong Initiator value. expect %s but get %s", redirects[1].Initiator, jsonRedirects[1].Initiator)
	}

	if jsonRedirects[1].Cookies[0].Raw != redirects[1].Cookies[0].Raw {
		t.Errorf("wrong Cookie[0] value. expect %s but get %s", redirects[1].Cookies[0].Raw, jsonRedirects[1].Cookies[0].Raw)
	}

	if jsonRedirects[1].Cookies[1].Raw != redirects[1].Cookies[1].Raw {
		t.Errorf("wrong Cookie[1] value. expect %s but get %s", redirects[1].Cookies[1].Raw, jsonRedirects[1].Cookies[1].Raw)
	}

}
