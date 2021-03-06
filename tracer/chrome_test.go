package tracer

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/raff/godet"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("/usr/bin/google-chrome", "--addr=localhost", "--port=9222", "--remote-debugging-port=9222", "--remote-debugging-address=0.0.0.0", "--disable-extensions", "--disable-gpu", "--headless", "--hide-scrollbars", "--no-first-run", "--no-sandbox")

	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("google-chrome headless runned with PID: %d\n", cmd.Process.Pid)
	log.Println("google-chrome headless runned on 9222 port")

	time.Sleep(1 * time.Second)

	code := m.Run()

	log.Printf("killing google-chrom PID %d\n", cmd.Process.Pid)
	// Kill chrome:
	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}

	os.Exit(code)
}

func TestNewChromeTracer(t *testing.T) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		t.Fatalf("cannot connect to Chrome instance: %s", err)
		return
	}

	size := &ScreenSize{
		Width:  1920,
		Height: 1080,
	}

	chr := NewChromeTracer(remote, size, "./assets")

	if chr.instance != remote {
		t.Error("wrong remote debuger instance")
	}

	err = remote.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestChromeTracer_Trace(t *testing.T) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		t.Fatalf("cannot connect to Chrome instance: %s", err)
		return
	}

	defer remote.Close()

	size := &ScreenSize{
		Width:  1920,
		Height: 1080,
	}

	chr := NewChromeTracer(remote, size, "./assets")

	if chr.instance != remote {
		t.Error("wrong remote debuger instance")
	}

	traceURL, err := url.Parse("https://www.google.com.ua")
	if err != nil {
		t.Error(err)
	}

	redirects, err := chr.Trace(traceURL, "testScreenshot.png")
	if err != nil {
		t.Error(err)
	}

	if len(redirects) != 0 {
		t.Error("No redirects expected")
	}
}

func TestChromeTracer_Trace2(t *testing.T) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		t.Fatalf("cannot connect to Chrome instance: %s", err)
		return
	}

	defer remote.Close()

	size := &ScreenSize{
		Width:  1920,
		Height: 1080,
	}

	chr := NewChromeTracer(remote, size, "./assets")

	if chr.instance != remote {
		t.Error("wrong remote debuger instance")
	}

	traceURL, err := url.Parse("http://google.com")
	if err != nil {
		t.Error(err)
	}

	redirects, err := chr.Trace(traceURL, "testScreenshot.png")
	if err != nil {
		t.Error(err)
	}

	if len(redirects) != 3 {
		t.Errorf("Two redirects expected but get %d", len(redirects))

		for _, redir := range redirects {
			t.Errorf("From %s -> To %s", redir.From.String(), redir.To.String())
		}
	}
}

func TestParseCookies(t *testing.T) {
	expectCookie := http.Cookie{
		Name:       "foo",
		Value:      "bar",
		Domain:     "test.com",
		MaxAge:     259200,
		RawExpires: "Mon, 31-Dec-2055 23:59:59 GMT",
		Path:       "/test",
	}
	rawCookie := "foo=bar; expires=Mon, 31-Dec-2055 23:59:59 GMT; Max-Age=259200; domain=test.com; Path=/test"
	cookies := parseCookies(rawCookie)

	if cookies[0].Value != expectCookie.Value {
		t.Errorf("invalid cookie Value. expect %s but get %s", expectCookie.Value, cookies[0].Value)
	}

	if cookies[0].Name != expectCookie.Name {
		t.Errorf("invalid cookie Name. expect %s but get %s", expectCookie.Name, cookies[0].Name)
	}

	if cookies[0].Domain != expectCookie.Domain {
		t.Errorf("invalid cookie Domain. expect %s but get %s", expectCookie.Domain, cookies[0].Domain)
	}

	if cookies[0].MaxAge != expectCookie.MaxAge {
		t.Errorf("invalid cookie MaxAge. expect %d but get %d", expectCookie.MaxAge, cookies[0].MaxAge)
	}

	if cookies[0].RawExpires != "Mon, 31-Dec-2055 23:59:59 GMT" {
		t.Errorf("invalid cookie RawExpires. expect %s but get %s", expectCookie.RawExpires, cookies[0].RawExpires)
	}

	if cookies[0].Raw != rawCookie {
		t.Errorf("invalid cookie Raw. expect %s but get %s", rawCookie, cookies[0].Raw)
	}

	if cookies[0].Path != expectCookie.Path {
		t.Errorf("invalid cookie Path. expect %s but get %s", expectCookie.Path, cookies[0].Path)
	}
}

func TestParseRedirectFromRaw1(t *testing.T) {
	redirectParams := godet.Params{}

	if err := json.Unmarshal([]byte(params1), &redirectParams); err != nil {
		panic(err)
	}

	redirect, err := parseRedirectFromRaw(redirectParams)
	if err != nil {
		t.Error(err)
	}

	if redirect.To.String() != "http://step1.test" {
		t.Errorf("invalid redirect To param. expect %s but get %s", "http://step1.test", redirect.To.String())
	}

	if redirect.From.String() != "http://step0.test" {
		t.Errorf("invalid redirect From param. expect %s but get %s", "http://step0.test", redirect.From.String())
	}

	if redirect.Status != 302 {
		t.Errorf("invalid redirect To param. expect %d but get %d", 302, redirect.Status)
	}

	if redirect.Initiator != "other" {
		t.Errorf("invalid redirect Initiator param. expect %s but get %s", "other", redirect.Initiator)
	}

	if redirect.RequestHeaders.Get("Test") != "redirective-request-header" {
		t.Errorf("invalid redirect RequestHeader param. expect %s but get %s", "redirective-request-header", redirect.RequestHeaders.Get("Test"))
	}

	if redirect.ResponseHeaders.Get("Test") != "redirective-response-header" {
		t.Errorf("invalid redirect ResponseHeaders param. expect %s but get %s", "redirective-response-header", redirect.ResponseHeaders.Get("Test"))
	}

	if len(redirect.Cookies) != 1 {
		t.Errorf("invalid redirect Cookies amount. expect %d but get %d", 1, len(redirect.Cookies))
	}

	if redirect.Cookies[0].Value != "bar" || redirect.Cookies[0].Name != "foo" || redirect.Cookies[0].Raw != "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com" {
		t.Error("invalid redirect Cookies values")
	}
}

func TestParseRedirectFromRaw2(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params2), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil {
		t.Errorf("expect error: %s", errorMessageInvalidToURL)
	}

	if err.Error() != errorMessageInvalidToURL {
		t.Errorf("expect error: %s but got %s", errorMessageInvalidToURL, err)
	}
}

func TestParseRedirectFromRaw3(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params3), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil {
		t.Errorf("expect error: %s", errorMessageInvalidFromURL)
	}

	if err.Error() != errorMessageInvalidFromURL {
		t.Errorf("expect error: %s but got %s", errorMessageInvalidFromURL, err)
	}
}

func TestParseRedirectFromRaw4(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params4), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil {
		t.Errorf("expect error: %s", errorMessageRedirectResponseNotExists)
	}

	if err.Error() != errorMessageRedirectResponseNotExists {
		t.Errorf("expect error: %s but got %s", errorMessageRedirectResponseNotExists, err)
	}
}

func TestParseRedirectFromRaw5(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params5), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil {
		t.Errorf("expect error: %s", errorMessageRequestNotExists)
	}

	if err.Error() != errorMessageRequestNotExists {
		t.Errorf("expect error: %s but got %s", errorMessageRequestNotExists, err)
	}
}

func TestParseRedirectFromRaw6(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params6), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil {
		t.Errorf("expect error: %s", errorMessageRedirectResponseParamHeadersNotExists)
	}

	if err.Error() != errorMessageRedirectResponseParamHeadersNotExists {
		t.Errorf("expect error: %s but got %s", errorMessageRedirectResponseParamHeadersNotExists, err)
	}
}
func TestParseRedirectFromRaw7(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params7), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil || err.Error() != errorMessageRedirectResponseParamURLNotExists {
		t.Errorf("expect error: %s", errorMessageRedirectResponseParamURLNotExists)
	}

	if err.Error() != errorMessageRedirectResponseParamURLNotExists {
		t.Errorf("expect error: %s but got %s", errorMessageRedirectResponseParamURLNotExists, err)
	}
}

func TestParseRedirectFromRaw8(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params8), &redirectParams); err != nil {
		panic(err)
	}

	_, err := parseRedirectFromRaw(redirectParams)

	if err == nil {
		t.Errorf("Expect error: %s", errorMessageRequestParamHeadersNotExists)
	}

	if err.Error() != errorMessageRequestParamHeadersNotExists {
		t.Errorf("Expect error: %s but got %s", errorMessageRequestParamHeadersNotExists, err)
	}
}

const params1 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headers": {
     "Test": "redirective-response-header",
     "Connection": "keep-alive",
     "Content-Type": "text/html; charset=UTF-8",
     "Date": "Fri, 28 Dec 2018 18:32:22 GMT",
     "Location": "http://step1.test",
     "Server": "nginx/1.10.3 (Ubuntu)",
     "Set-Cookie": "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com",
     "Transfer-Encoding": "chunked"
   },
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   },
   "url": "http://step0.test"
 },
 "request": {
   "headers": {
     "Test": "redirective-request-header",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "http://step0.test"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

const params2 = `
{
 "documentURL": "user@httpstep1test:ew?something#wron=here",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headers": {
     "Test": "redirective-response-header",
     "Connection": "keep-alive",
     "Content-Type": "text/html; charset=UTF-8",
     "Date": "Fri, 28 Dec 2018 18:32:22 GMT",
     "Location": "http://step1.test",
     "Server": "nginx/1.10.3 (Ubuntu)",
     "Set-Cookie": "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com",
     "Transfer-Encoding": "chunked"
   },
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   },
   "url": "http://step0.test"
 },
 "request": {
   "headers": {
     "Test": "redirective-request-header",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "http://step0.test"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

const params3 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headers": {
     "Test": "redirective-response-header",
     "Connection": "keep-alive",
     "Content-Type": "text/html; charset=UTF-8",
     "Date": "Fri, 28 Dec 2018 18:32:22 GMT",
     "Location": "http://step1.test",
     "Server": "nginx/1.10.3 (Ubuntu)",
     "Set-Cookie": "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com",
     "Transfer-Encoding": "chunked"
   },
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   },
   "url": "user@httpstep1test:ew?something#wron=here"
 },
 "request": {
   "headers": {
     "Test": "redirective-request-header",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "user@httpstep1test:ew?something#wron=here"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

const params4 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "request": {
   "headers": {
     "Test": "redirective-request-header",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "user@httpstep1test:ew?something#wron=here"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

const params5 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headers": {
     "Test": "redirective-response-header",
     "Connection": "keep-alive",
     "Content-Type": "text/html; charset=UTF-8",
     "Date": "Fri, 28 Dec 2018 18:32:22 GMT",
     "Location": "http://step1.test",
     "Server": "nginx/1.10.3 (Ubuntu)",
     "Set-Cookie": "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com",
     "Transfer-Encoding": "chunked"
   },
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   },
   "url": "user@httpstep1test:ew?something#wron=here"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

const params6 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   },
   "url": "http://step0.test"
 },
 "request": {
   "headers": {
     "Test": "redirective-request-header",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "http://step0.test"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`
const params7 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headers": {
     "Test": "redirective-response-header",
     "Connection": "keep-alive",
     "Content-Type": "text/html; charset=UTF-8",
     "Date": "Fri, 28 Dec 2018 18:32:22 GMT",
     "Location": "http://step1.test",
     "Server": "nginx/1.10.3 (Ubuntu)",
     "Set-Cookie": "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com",
     "Transfer-Encoding": "chunked"
   },
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   }
 },
 "request": {
   "headers": {
     "Test": "redirective-request-header",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "http://step0.test"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

const params8 = `
{
 "documentURL": "http://step1.test",
 "frameId": "F394EA807250832376BE81745B17B0E9",
 "hasUserGesture": false,
 "initiator": {
   "type": "other"
 },
 "loaderId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "redirectResponse": {
   "connectionId": 56,
   "connectionReused": false,
   "encodedDataLength": 427,
   "fromDiskCache": false,
   "fromServiceWorker": false,
   "headers": {
     "Test": "redirective-response-header",
     "Connection": "keep-alive",
     "Content-Type": "text/html; charset=UTF-8",
     "Date": "Fri, 28 Dec 2018 18:32:22 GMT",
     "Location": "http://step1.test",
     "Server": "nginx/1.10.3 (Ubuntu)",
     "Set-Cookie": "foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com",
     "Transfer-Encoding": "chunked"
   },
   "headersText": "HTTP/1.1 302 Found\r\nServer: nginx/1.10.3 (Ubuntu)\r\nDate: Fri, 28 Dec 2018 18:32:22 GMT\r\nContent-Type: text/html; charset=UTF-8\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\nSet-Cookie: foo=bar; expires=Sat, 28-Dec-2019 18:32:22 GMT; Max-Age=31536000; domain=test.com\r\nLocation: http://step1.com\r\n\r\n",
   "mimeType": "text/html",
   "protocol": "http/1.1",
   "remoteIPAddress": "104.248.96.70",
   "remotePort": 80,
   "requestHeaders": {
     "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
     "Accept-Encoding": "gzip, deflate",
     "Connection": "keep-alive",
     "Host": "step0.com",
     "Upgrade-Insecure-Requests": "1",
     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36"
   },
   "requestHeadersText": "GET / HTTP/1.1\r\nHost: step0.test\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate\r\n",
   "securityState": "neutral",
   "status": 302,
   "statusText": "Found",
   "timing": {
     "connectEnd": 464.923,
     "connectStart": 113.478,
     "dnsEnd": 113.478,
     "dnsStart": 0.105,
     "proxyEnd": -1,
     "proxyStart": -1,
     "pushEnd": 0,
     "pushStart": 0,
     "receiveHeadersEnd": 1009.774,
     "requestTime": 15875.097865,
     "sendEnd": 465.061,
     "sendStart": 464.998,
     "sslEnd": -1,
     "sslStart": -1,
     "workerReady": -1,
     "workerStart": -1
   },
   "url": "http://step0.test"
 },
 "request": {
   "initialPriority": "VeryHigh",
   "method": "GET",
   "mixedContentType": "none",
   "referrerPolicy": "no-referrer-when-downgrade",
   "url": "http://step1.test"
 },
 "requestId": "E8DAACD689A021E0963DA6DDC3FC9AF9",
 "timestamp": 15876.109173,
 "type": "Document",
 "wallTime": 1546021943.764962
}`

func Test_pareseMainResponseFromRaw_NoResponse(t *testing.T) {
	var input godet.Params = godet.Params{}

	rd, err := pareseMainResponseFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageResponseParamNotExists)
	} else if strings.Compare(err.Error(), errorMessageResponseParamNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageResponseParamNotExists, err.Error())
	}

	if rd != nil {
		t.Error("not expected redirect")
	}
}

func Test_pareseMainResponseFromRaw_Response_NoUrl(t *testing.T) {
	var input godet.Params = godet.Params{
		"response": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookies": "set-cookie",
				"someName":    "someValue",
			},
			//"url": "http://www.google.com.ua",
			"requestHeaders": "",
			"status":         "200",
		},
	}

	rd, err := pareseMainResponseFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRedirectResponseParamURLNotExists)
	} else if strings.Compare(err.Error(), errorMessageRedirectResponseParamURLNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRedirectResponseParamURLNotExists, err.Error())
	}

	if rd != nil {
		t.Error("not expected redirect")
	}
}

func Test_pareseMainResponseFromRaw_Response_InvalidUrl(t *testing.T) {
	var input godet.Params = godet.Params{
		"response": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookies": "set-cookie",
				"someName":    "someValue",
			},
			"url":            ":wwwgooglecomua",
			"requestHeaders": "",
			"status":         "200",
		},
	}

	rd, err := pareseMainResponseFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageInvalidToURL)
	} else if strings.Compare(err.Error(), errorMessageInvalidToURL) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageInvalidToURL, err.Error())
	}

	if rd != nil {
		t.Error("not expected redirect")
	}
}

func Test_pareseMainResponseFromRaw_Response_NoHeaders(t *testing.T) {
	var input godet.Params = godet.Params{
		"response": map[string]interface{}{
			//"headers": map[string]interface{}{
			//	"set-cookies": "set-cookie",
			//	"someName": "someValue",
			//},
			"url":            "http://www.google.com.ua",
			"requestHeaders": "",
			"status":         "200",
		},
	}

	rd, err := pareseMainResponseFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRedirectResponseParamHeadersNotExists)
	} else if strings.Compare(err.Error(), errorMessageRedirectResponseParamHeadersNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRedirectResponseParamHeadersNotExists, err.Error())
	}

	if rd != nil {
		t.Error("not expected redirect")
	}
}

func Test_pareseMainResponseFromRaw_Response_NoRequestHeaders(t *testing.T) {
	var input godet.Params = godet.Params{
		"response": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookies": "set-cookie",
				"someName":    "someValue",
			},
			"url": "http://www.google.com.ua",
			//"requestHeaders": "",
			"status": "200",
		},
	}

	rd, err := pareseMainResponseFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRequestParamHeadersNotExists)
	} else if strings.Compare(err.Error(), errorMessageRequestParamHeadersNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRequestParamHeadersNotExists, err.Error())
	}

	if rd != nil {
		t.Error("not expected redirect")
	}
}

func Test_pareseMainResponseFromRaw(t *testing.T) {
	input := godet.Params{
		"response": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key0=f66d4763-7f3f-4ac1-b30a-3cbf31f46123; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
			"url": "http://www.google.com.ua",
			"requestHeaders": map[string]interface{}{
				"someName1": "someValue1",
				"someName2": "someValue2",
				"someName3": "someValue3",
			},
			"status": 203.00,
		},
	}

	rd, err := pareseMainResponseFromRaw(input)
	if err != nil {
		t.Errorf("unexpected error `%s`", err)
	}

	if rd.Status != 203 {
		t.Errorf("wrong redirect status parsed")
	}

	if len(rd.Cookies) != 1 {
		t.Errorf("wrong cookies count. expected 1 cookie")
	}

	if ck := rd.Cookies[0]; ck.Name != "key0" {
		t.Errorf("wrong cookie name, expected `key0`, but got `%s`", ck.Name)
	}

	if ck := rd.Cookies[0]; ck.Value != "f66d4763-7f3f-4ac1-b30a-3cbf31f46123" {
		t.Errorf("wrong cookie value, expected `f66d4763-7f3f-4ac1-b30a-3cbf31f46123`, but got `%s`", ck.Value)
	}

	if rd.To.String() != "http://www.google.com.ua" {
		t.Errorf("wrong from url. expected `http://www.google.com.ua`, but got `%s`", rd.To.String())
	}

	reqHeaders := *rd.RequestHeaders
	if len(reqHeaders) != 3 {
		t.Errorf("expected to get 2 request headers, but got %d", len(reqHeaders))
	}

	if reqHeaders.Get("someName1") != "someValue1" {
		t.Errorf("wrong value in `someName1`. expected `someValue1`, but got `%s`", reqHeaders.Get("someName1"))
	}

	if reqHeaders.Get("someName2") != "someValue2" {
		t.Errorf("wrong value in `someName2`. expected `someValue2`, but got `%s`", reqHeaders.Get("someName2"))
	}

	if reqHeaders.Get("someName3") != "someValue3" {
		t.Errorf("wrong value in `someName3`. expected `someValue3`, but got `%s`", reqHeaders.Get("someName3"))
	}

	respHeaders := *rd.ResponseHeaders
	if len(respHeaders) != 2 {
		t.Errorf("expected to get 2 request headers, but got %d", len(respHeaders))
	}

	if respHeaders.Get("set-cookie") != "key0=f66d4763-7f3f-4ac1-b30a-3cbf31f46123; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com" {
		t.Errorf("wrong value in `set-cookie`. expected `key0=f66d4763-7f3f-4ac1-b30a-3cbf31f46123; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com`, but got `%s`", respHeaders.Get("set-cookie"))
	}

	if respHeaders.Get("someName") != "someValue" {
		t.Errorf("wrong value in `someName`. expected `someValue`, but got `%s`", respHeaders.Get("someName"))
	}
}

func Test_ParseCookies(t *testing.T) {
	rawCookie := "key1=f66d4763-7f3f-4ac1-b30a-3cbf31f46200; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; path=/; domain=domain.com"

	cookies := parseCookies(rawCookie)
	if len(cookies) != 1 {
		t.Errorf("wrong cookies count parsed. Expected 1 cookie, but got %d", len(cookies))
	}

	if cookies[0].Name != "key1" {
		t.Errorf("wrong cookie name parsed. expected name `key1`, but got `%s`", cookies[0].Name)
	}

	if cookies[0].Value != "f66d4763-7f3f-4ac1-b30a-3cbf31f46200" {
		t.Errorf("wrong cookie value parsed. expected value `f66d4763-7f3f-4ac1-b30a-3cbf31f46200`, but got `%s`", cookies[0].Value)
	}

	if cookies[0].Path != "/" {
		t.Errorf("wrong cookie path parsed. expected value `/`, but got `%s`", cookies[0].Path)
	}

	if cookies[0].MaxAge != 2592000 {
		t.Errorf("wrong cookie max age parsed. expected value `2592000`, but got `%d`", cookies[0].MaxAge)
	}

	if cookies[0].RawExpires != "Wed, 08-Jan-2100 18:01:07 GMT" {
		t.Errorf("wrong cookie expires parsed. expected value `Wed, 08-Jan-2100 18:01:07 GMT`, but got `%s`", cookies[0].RawExpires)
	}

	if cookies[0].Domain != "domain.com" {
		t.Errorf("wrong cookie domain parsed. expected value `domain.com`, but got `%s`", cookies[0].Domain)
	}

	year, month, day := cookies[0].Expires.Date()
	if year != 2100 {
		t.Errorf("wrong cookie expire year parsed. expected value `2100`, but got `%d`", year)
	}

	if day != 8 {
		t.Errorf("wrong cookie expire day parsed. expected value `8`, but got `%d`", 8)
	}

	if month.String() != "January" {
		t.Errorf("wrong cookie expire month parsed. expected value `January`, but got `%s`", month.String())
	}

	if cookies[0].Expires.Hour() != 18 {
		t.Errorf("wrong cookie expire hours parsed. expected value `18`, but got `%d`", cookies[0].Expires.Hour())
	}

	if cookies[0].Expires.Minute() != 1 {
		t.Errorf("wrong cookie expire minutes parsed. expected value `1`, but got `%d`", cookies[0].Expires.Minute())
	}

	if cookies[0].Expires.Second() != 7 {
		t.Errorf("wrong cookie expire seconds parsed. expected value `7`, but got `%d`", cookies[0].Expires.Second())
	}
}

func Test_ParseHeadersFromRaw_NoHeaders(t *testing.T) {
	request := map[string]interface{}{
		"not-headers": []string{"some value", "some value2"},
	}

	result, err := parseHeadersFromRaw(request)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRequestParamHeadersNotExists)
	} else if strings.Compare(errorMessageRequestParamHeadersNotExists, err.Error()) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRequestParamHeadersNotExists, err.Error())
	}

	if len(*result) != 0 {
		t.Errorf("not expected headers count")
	}
}

func Test_ParseHeadersFromRaw(t *testing.T) {
	request := map[string]interface{}{
		"headers": map[string]interface{}{
			"header1": "some value1",
			"header2": "some value2",
		},
	}

	result, err := parseHeadersFromRaw(request)
	if err != nil {
		t.Errorf("unexpected error `%s`", err.Error())
	}

	if len(*result) != 2 {
		t.Errorf("expected 2 headers, but got `%d`", len(*result))
	}
}

func Test_ParseRedirectFromRaw_NoRedirectResponse(t *testing.T) {
	var input godet.Params = godet.Params{
		"response": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628a; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
			"url": "http://www.google.com.ua",
			"requestHeaders": map[string]interface{}{
				"someName1": "someValue1",
				"someName2": "someValue2",
				"someName3": "someValue3",
			},
			"status": 203.00,
		},
	}

	r, err := parseRedirectFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRedirectResponseNotExists)
	} else if strings.Compare(err.Error(), errorMessageRedirectResponseNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRedirectResponseNotExists, err.Error())
	}

	if r != nil {
		t.Error("no redirects expected")
	}
}

func Test_ParseRedirectFromRaw_NoRequest(t *testing.T) {
	var input godet.Params = godet.Params{
		"redirectResponse": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628b; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
			"url":    "http://www.google.com.ua",
			"status": 203.00,
		},
	}

	r, err := parseRedirectFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRequestNotExists)
	} else if strings.Compare(err.Error(), errorMessageRequestNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRequestNotExists, err.Error())
	}

	if r != nil {
		t.Error("no redirects expected")
	}
}

func Test_ParseRedirectFromRaw_InvalidDocumentURL(t *testing.T) {
	var input godet.Params = godet.Params{
		"redirectResponse": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628c; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
			"url":    "http://www.google.com.ua",
			"status": 203.00,
		},
		"documentURL": ":www.google.com.ua",
		"request": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628d; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
		},
	}

	r, err := parseRedirectFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageInvalidToURL)
	} else if strings.Compare(err.Error(), errorMessageInvalidToURL) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageInvalidToURL, err.Error())
	}

	if r != nil {
		t.Error("no redirects expected")
	}
}

func Test_ParseRedirectFromRaw_RedirectNoURL(t *testing.T) {
	var input godet.Params = godet.Params{
		"redirectResponse": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628e; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
			//"url": "http://www.google.com.ua",
			"status": 203.00,
		},
		"documentURL": "http://www.google.com.ua",
		"request": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628f; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
		},
	}

	r, err := parseRedirectFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRedirectResponseParamURLNotExists)
	} else if strings.Compare(err.Error(), errorMessageRedirectResponseParamURLNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRedirectResponseParamURLNotExists, err.Error())
	}

	if r != nil {
		t.Error("no redirects expected")
	}
}

func Test_ParseRedirectFromRaw_RedirectInvalidURL(t *testing.T) {
	var input godet.Params = godet.Params{
		"redirectResponse": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628g; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
			"url":    ":www.google.com.ua",
			"status": 203.00,
		},
		"documentURL": "http://www.google.com.ua",
		"request": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628h; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
		},
	}

	r, err := parseRedirectFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageInvalidFromURL)
	} else if strings.Compare(err.Error(), errorMessageInvalidFromURL) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageInvalidFromURL, err.Error())
	}

	if r != nil {
		t.Error("no redirects expected")
	}
}

func Test_ParseRedirectFromRaw_RedirectNoHeaders(t *testing.T) {
	var input godet.Params = godet.Params{
		"redirectResponse": map[string]interface{}{
			"url":    "http://www.google.com.ua",
			"status": 203.00,
		},
		"documentURL": "http://www.google.com.ua",
		"request": map[string]interface{}{
			"headers": map[string]interface{}{
				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628c; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
				"someName":   "someValue",
			},
		},
	}

	r, err := parseRedirectFromRaw(input)
	if err == nil {
		t.Errorf("expected error `%s`", errorMessageRedirectResponseParamHeadersNotExists)
	} else if strings.Compare(err.Error(), errorMessageRedirectResponseParamHeadersNotExists) != 0 {
		t.Errorf("expected error `%s`, but got `%s`", errorMessageRedirectResponseParamHeadersNotExists, err.Error())
	}

	if r != nil {
		t.Error("no redirects expected")
	}
}

//func Test_ParseRedirectFromRaw(t *testing.T) {
//	var input godet.Params = godet.Params{
//		"redirectResponse": map[string]interface{}{
//			"headers": map[string]interface{}{
//				"set-cookie": "key=00000000-0000-4ac1-b30a-3cbf31f4628c; expires=Wed, 13-Oct-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
//				"someName":   "someValue",
//			},
//			"url":    "http://www.google.com.ua/from",
//			"status": 203.00,
//		},
//		"documentURL": "http://www.google.com.ua/to",
//		"request": map[string]interface{}{
//			"headers": map[string]interface{}{
//				"set-cookie": "key=f66d4763-7f3f-4ac1-b30a-3cbf31f4628c; expires=Wed, 08-Jan-2100 18:01:07 GMT; Max-Age=2592000; domain=domain.com",
//				"someName":   "someValue",
//			},
//		},
//	}
//
//	rd, err := parseRedirectFromRaw(input)
//	if err != nil {
//		t.Errorf("unexpected error `%s`", err.Error())
//	}
//	//TODO: check rd!
//}
