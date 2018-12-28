package tracer

import (
	"testing"
)

func TestChromeTracer_GetTrace_ParseCookies(t *testing.T) {
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
		t.Errorf("invalid cookie RawExpires. expect %d but get %d", "Mon, 31-Dec-2055 23:59:59 GMT", cookies[0].RawExpires)
	}
	if cookies[0].Raw != rawCookie {
		t.Errorf("invalid cookie Raw. expect %s but get %s", rawCookie, cookies[0].Raw)
	}
	if cookies[0].Path != "/test" {
		t.Errorf("invalid cookie Path. expect %s but get %s", "/test", cookies[0].Path)
	}
	//if cookies[0].Expires != time.Parse()
	//TODO: Expires

}