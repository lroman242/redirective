package tracer

import (
	"encoding/json"
	"github.com/raff/godet"
	"net/url"
	"testing"
)

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
	traceUrl, err := url.Parse("https://www.google.com.ua")
	if err != nil {
		t.Error(err)
	}

	redirects, err := chr.GetTrace(traceUrl)
	if err == nil || err.Error() != "No redirects found" {
		t.Errorf("Expect error: No redirects found")
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

	traceUrl, err := url.Parse("http://google.com")
	if err != nil {
		t.Error(err)
	}

	redirects, err := chr.GetTrace(traceUrl)
	if err != nil {
		t.Error(err)
	}

	if len(redirects) != 2 {
		t.Errorf("Two redirects expected but get %d", len(redirects))
		for _, redir := range redirects {
			t.Errorf("From %s -> To %s", redir.From.String(), redir.To.String())
		}
	}
}

func TestParseCookies(t *testing.T) {
	rawCookie := "foo=bar; expires=Mon, 31-Dec-2055 23:59:59 GMT; Max-Age=259200; domain=test.com; Path=/test"
	cookies := parseCookies(rawCookie)

	if cookies[0].Value != "bar" {
		t.Errorf("invalid cookie Value. expect %s but get %s", "bar", cookies[0].Value)
	}
	if cookies[0].Name != "foo" {
		t.Errorf("invalid cookie Name. expect %s but get %s", "foo", cookies[0].Name)
	}
	if cookies[0].Domain != "test.com" {
		t.Errorf("invalid cookie Domain. expect %s but get %s", "test.com", cookies[0].Domain)
	}
	if cookies[0].MaxAge != 259200 {
		t.Errorf("invalid cookie MaxAge. expect %d but get %d", 259200, cookies[0].MaxAge)
	}
	if cookies[0].RawExpires != "Mon, 31-Dec-2055 23:59:59 GMT" {
		t.Errorf("invalid cookie RawExpires. expect %s but get %s", "Mon, 31-Dec-2055 23:59:59 GMT", cookies[0].RawExpires)
	}
	if cookies[0].Raw != rawCookie {
		t.Errorf("invalid cookie Raw. expect %s but get %s", rawCookie, cookies[0].Raw)
	}
	if cookies[0].Path != "/test" {
		t.Errorf("invalid cookie Path. expect %s but get %s", "/test", cookies[0].Path)
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
	if err == nil || err.Error() != "Invalid redirect To url" {
		t.Errorf("Expect error: Invalid redirect To url")
	}
}

func TestParseRedirectFromRaw3(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params3), &redirectParams); err != nil {
		panic(err)
	}
	_, err := parseRedirectFromRaw(redirectParams)
	if err == nil || err.Error() != "Invalid redirect From url" {
		t.Errorf("Expect error: Invalid redirect From url")
	}
}

func TestParseRedirectFromRaw4(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params4), &redirectParams); err != nil {
		panic(err)
	}
	_, err := parseRedirectFromRaw(redirectParams)
	if err == nil || err.Error() != "Invalid redirect. redirectResponse param not exists" {
		t.Errorf("Expect error: Invalid redirect. redirectResponse param not exists")
	}
}

func TestParseRedirectFromRaw5(t *testing.T) {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params5), &redirectParams); err != nil {
		panic(err)
	}
	_, err := parseRedirectFromRaw(redirectParams)
	if err == nil || err.Error() != "Invalid redirect. request param not exists" {
		t.Errorf("Expect error: Invalid redirect. request param not exists")
	}
}

func TestParseRedirectFromRaw6(t *testing.T)  {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params6), &redirectParams); err != nil {
		panic(err)
	}
	_, err := parseRedirectFromRaw(redirectParams)
	if err == nil || err.Error() != "Invalid redirect. redirectResponse param headers not exists" {
		t.Errorf("Invalid redirect. redirectResponse param headers not exists")
	}
}
func TestParseRedirectFromRaw7(t *testing.T)  {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params7), &redirectParams); err != nil {
		panic(err)
	}
	_, err := parseRedirectFromRaw(redirectParams)
	if err == nil || err.Error() != "Invalid redirect. redirectResponse param url not exists" {
		t.Errorf("Expect error: Invalid redirect. redirectResponse param url not exists")
	}
}

func TestParseRedirectFromRaw8(t *testing.T)  {
	redirectParams := godet.Params{}
	if err := json.Unmarshal([]byte(params8), &redirectParams); err != nil {
		panic(err)
	}
	_, err := parseRedirectFromRaw(redirectParams)
	if err == nil || err.Error() != "Invalid redirect. request param headers not exists" {
		t.Errorf("Expect error: Invalid redirect. request param headers not exists")
	}
}

const params1  = `
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

const params2  = `
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

const params3  = `
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

const params4  = `
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

const params5  = `
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

const params6  = `
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
const params7  = `
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

const params8  = `
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
