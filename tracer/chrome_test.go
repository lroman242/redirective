package tracer

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

	chr := NewChromeTracer(remote)

	if chr.instance != remote {
		t.Error("wrong remote debuger instance")
	}

	err = remote.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestChromeTracer_GetTrace(t *testing.T) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		t.Fatalf("cannot connect to Chrome instance: %s", err)
		return
	}

	defer remote.Close()

	chr := NewChromeTracer(remote)

	if chr.instance != remote {
		t.Error("wrong remote debuger instance")
	}

	traceURL, err := url.Parse("https://www.google.com.ua")
	if err != nil {
		t.Error(err)
	}

	redirects, err := chr.GetTrace(traceURL)
	if err != nil {
		t.Error(err)
	}

	if len(redirects) != 0 {
		t.Error("No redirects expected")
	}
}

func TestChromeTracer_GetTrace2(t *testing.T) {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		t.Fatalf("cannot connect to Chrome instance: %s", err)
		return
	}

	defer remote.Close()

	chr := NewChromeTracer(remote)

	if chr.instance != remote {
		t.Error("wrong remote debuger instance")
	}

	traceURL, err := url.Parse("http://google.com")
	if err != nil {
		t.Error(err)
	}

	redirects, err := chr.GetTrace(traceURL)
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
