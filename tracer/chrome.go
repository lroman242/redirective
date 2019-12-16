package tracer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/log"
	"github.com/raff/godet"
)

const setCookieHeaderName = "set-cookie"
const documentParamName = "Document"

const (
	errorMessageInvalidMainFrameID                      = "invalid mainframe id"
	errorMessageNoResponseFromMainFrame                 = "no responses found for mainframe"
	errorMessageRedirectResponseNotExists               = "invalid redirect. `redirectResponse` param not exists"
	errorMessageRequestNotExists                        = "invalid redirect. `request` param not exists"
	errorMessageInvalidToURL                            = "invalid redirect `To` url"
	errorMessageInvalidFromURL                          = "invalid redirect `From` url"
	errorMessageRedirectResponseParamURLNotExists       = "invalid redirect. `redirectResponse` param `url` not exists"
	errorMessageRedirectResponseParamHeadersNotExists   = "invalid redirect. `redirectResponse` param `headers` not exists"
	errorMessageRequestParamHeadersNotExists            = "invalid redirect. request param `headers` not exists"
	errorMessageResponseParamNotExists                  = "invalid redirect. `response` param not exists"
	errorMessageRedirectResponseParamInitiatorNotExists = "invalid redirect. `initiator` param not exists"
)

// ChromeRemoteDebuggerInterface implements API to work with browser debugger
type ChromeRemoteDebuggerInterface interface {
	EnableRequestInterception(enabled bool) error
	CallbackEvent(method string, cb godet.EventCallback)
	NewTab(url string) (*godet.Tab, error)
	CloseTab(tab *godet.Tab) error
	NetworkEvents(enable bool) error
	ActivateTab(tab *godet.Tab) error
	AllEvents(enable bool) error
	Navigate(url string) (string, error)
	SetDeviceMetricsOverride(width int, height int, deviceScaleFactor float64, mobile bool, fitWindow bool) error
	SetVisibleSize(width, height int) error
	SaveScreenshot(filename string, perm os.FileMode, quality int, fromSurface bool) error
}

// ChromeTracer represent tracer based on google chrome debugging tools
type ChromeTracer struct {
	instance ChromeRemoteDebuggerInterface
}

// NewChromeTracer create new chrome tracer instance
func NewChromeTracer(chrome *godet.RemoteDebugger) *ChromeTracer {
	return &ChromeTracer{
		instance: chrome,
	}
}

func (ct *ChromeTracer) traceURL(url *url.URL, redirects, responses *map[string][]godet.Params) (string, error) {
	frameID := ""

	err := ct.instance.EnableRequestInterception(true)
	if err != nil {
		return frameID, fmt.Errorf("`EnableRequestInterception` failed. %s", err)
	}

	ct.instance.CallbackEvent("Network.requestWillBeSent", func(params godet.Params) {
		if _, ok := params["redirectResponse"]; ok && params["type"] == documentParamName {
			(*redirects)[params["frameId"].(string)] = append((*redirects)[params["frameId"].(string)], params)
		}
	})
	ct.instance.CallbackEvent("Network.responseReceived", func(params godet.Params) {
		if params["type"] == documentParamName {
			(*responses)[params["frameId"].(string)] = append((*responses)[params["frameId"].(string)], params)
		}
	})

	// create new tab
	tab, _ := ct.instance.NewTab("https://www.google.com")
	defer func(tab *godet.Tab) {
		err = ct.instance.CloseTab(tab)
		if err != nil {
			log.Error(fmt.Errorf("`CloseTab` failed. %s", err))
		}
	}(tab)

	err = ct.instance.NetworkEvents(true)
	if err != nil {
		return frameID, fmt.Errorf("`NetworkEvents failed. %s", err)
	}

	// navigate in existing tab
	err = ct.instance.ActivateTab(tab)
	if err != nil {
		return frameID, fmt.Errorf("`ActivateTab` failed. %s", err)
	}

	// re-enable events when changing active tab
	err = ct.instance.AllEvents(true) // enable all events
	if err != nil {
		return frameID, fmt.Errorf("`AllEvents` failed. %s", err)
	}

	frameID, err = ct.instance.Navigate(url.String())
	if err != nil {
		return frameID, fmt.Errorf("`Navigate` failed. %s", err)
	}

	return frameID, nil
}

// Trace parse redirect trace path for provided url
func (ct *ChromeTracer) Trace(url *url.URL) ([]*Redirect, error) {
	var redirects []*Redirect

	rawRedirects := make(map[string][]godet.Params)
	rawResponses := make(map[string][]godet.Params)

	frameID, err := ct.traceURL(url, &rawRedirects, &rawResponses)
	if err != nil {
		return redirects, err
	}

	if frameID == "" {
		return redirects, errors.New(errorMessageInvalidMainFrameID)
	}

	if len(rawRedirects) == 0 {
		return redirects, nil
	}

	if rawRedirects, ok := rawRedirects[frameID]; ok {
		for _, rawRedirect := range rawRedirects {
			redirect, err := parseRedirectFromRaw(rawRedirect)
			if err != nil {
				return redirects, fmt.Errorf("an error during parsing redirects. %s", err)
			}

			redirects = append(redirects, redirect)
		}
	} /* else {
		return redirects, errors.New("No redirects found for mainframe")
	}*/

	if rawRespons, ok := rawResponses[frameID]; ok {
		response, err := pareseMainResponseFromRaw(rawRespons[len(rawRespons)-1])
		if err != nil {
			return redirects, fmt.Errorf("an error during parsing response. %s", err)
		}

		redirects = append(redirects, response)
	} else {
		return redirects, errors.New(errorMessageNoResponseFromMainFrame)
	}

	return redirects, nil
}

// Screenshot function makes a final page screen capture
func (ct *ChromeTracer) Screenshot(url *url.URL, size *ScreenSize, path string) error {
	err := ct.instance.EnableRequestInterception(true)
	if err != nil {
		return fmt.Errorf("`EnableRequestInterception` failed. %s", err)
	}

	// create new tab
	tab, _ := ct.instance.NewTab(url.String())
	defer func(tab *godet.Tab) {
		err = ct.instance.CloseTab(tab)
		if err != nil {
			log.Error(fmt.Errorf("`CloseTab` failed. %s", err))
		}
	}(tab)

	// navigate in existing tab
	err = ct.instance.ActivateTab(tab)
	if err != nil {
		return fmt.Errorf("`ActivateTab` failed. %s", err)
	}

	err = ct.instance.SetDeviceMetricsOverride(size.Width, size.Height, 0, false, false)
	if err != nil {
		return fmt.Errorf("set screen size error: %s", err)
	}

	err = ct.instance.SetVisibleSize(size.Width, size.Height)
	if err != nil {
		return fmt.Errorf("set visibility size error: %s", err)
	}

	_, err = ct.instance.Navigate(url.String())
	if err != nil {
		return fmt.Errorf("`Navigate` failed. %s", err)
	}

	time.Sleep(time.Second * 5)

	// take a screenshot
	err = ct.instance.SaveScreenshot(path, 0644, 100, true)
	if err != nil {
		return fmt.Errorf("cannot capture screenshot: %s", err)
	}
	//time.Sleep(time.Second)

	return nil
}

func parseRedirectFromRaw(rawRedirect godet.Params) (*Redirect, error) {
	if _, ok := rawRedirect["redirectResponse"]; !ok {
		return nil, errors.New(errorMessageRedirectResponseNotExists)
	}

	if _, ok := rawRedirect["request"]; !ok {
		return nil, errors.New(errorMessageRequestNotExists)
	}

	redirectResponse := rawRedirect.Map("redirectResponse")
	request := rawRedirect.Map("request")

	to, err := url.Parse(rawRedirect.String("documentURL"))
	if err != nil {
		return nil, errors.New(errorMessageInvalidToURL)
	}

	if _, ok := redirectResponse["url"]; !ok {
		return nil, errors.New(errorMessageRedirectResponseParamURLNotExists)
	}

	from, err := url.Parse(redirectResponse["url"].(string))
	if err != nil {
		return nil, errors.New(errorMessageInvalidFromURL)
	}

	if _, ok := redirectResponse["headers"]; !ok {
		return nil, errors.New(errorMessageRedirectResponseParamHeadersNotExists)
	}

	var cookies []*http.Cookie

	responseHeaders := http.Header{}

	for index, header := range redirectResponse["headers"].(map[string]interface{}) {
		responseHeaders.Add(index, header.(string))

		if strings.ToLower(index) == setCookieHeaderName {
			cookies = parseCookies(header.(string))
		}
	}

	requestHeaders, err := parseHeadersFromRaw(request)
	if err != nil {
		return nil, err
	}

	status := int(redirectResponse["status"].(float64))

	if _, ok := rawRedirect["initiator"]; !ok {
		return nil, errors.New(errorMessageRedirectResponseParamInitiatorNotExists)
	}

	initiator := rawRedirect.Map("initiator")["type"].(string)

	return NewRedirect(from, to, requestHeaders, &responseHeaders, cookies, status, initiator), nil
}

func parseHeadersFromRaw(request map[string]interface{}) (*http.Header, error) {
	requestHeaders := &http.Header{}

	if _, ok := request["headers"]; !ok {
		return requestHeaders, errors.New(errorMessageRequestParamHeadersNotExists)
	}

	requestHeadersRaw := request["headers"].(map[string]interface{})

	for index, header := range requestHeadersRaw {
		requestHeaders.Add(index, header.(string))
	}

	return requestHeaders, nil
}

func parseCookies(s string) []*http.Cookie {
	rawCookies := strings.Split(s, "\n")
	cookies := make([]*http.Cookie, 0, len(rawCookies))

	for _, rawCookie := range rawCookies {
		parsedCookies := (&http.Response{Header: http.Header{"Set-Cookie": {rawCookie}}}).Cookies()
		cookies = append(cookies, parsedCookies...)
	}

	return cookies
}

func pareseMainResponseFromRaw(rawResponses godet.Params) (*Redirect, error) {
	if _, ok := rawResponses["response"]; !ok {
		return nil, errors.New(errorMessageResponseParamNotExists)
	}

	response := rawResponses.Map("response")

	if _, ok := response["url"]; !ok {
		return nil, errors.New(errorMessageRedirectResponseParamURLNotExists)
	}

	to, err := url.Parse(response["url"].(string))
	if err != nil {
		return nil, errors.New(errorMessageInvalidToURL)
	}

	if _, ok := response["headers"]; !ok {
		return nil, errors.New(errorMessageRedirectResponseParamHeadersNotExists)
	}

	var cookies []*http.Cookie

	responseHeaders := http.Header{}

	for index, header := range response["headers"].(map[string]interface{}) {
		responseHeaders.Add(index, header.(string))

		if strings.ToLower(index) == setCookieHeaderName {
			cookies = parseCookies(header.(string))
		}
	}

	if _, ok := response["requestHeaders"]; !ok {
		return nil, errors.New(errorMessageRequestParamHeadersNotExists)
	}

	requestHeaders := http.Header{}

	for index, header := range response["requestHeaders"].(map[string]interface{}) {
		requestHeaders.Add(index, header.(string))
	}

	status := int(response["status"].(float64))

	redirect := NewRedirect(&url.URL{}, to, &requestHeaders, &responseHeaders, cookies, status, "")

	return redirect, nil
}
