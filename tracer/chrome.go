package tracer

import (
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/raff/godet"
	"net/http"
	"net/url"
	"time"
)

type chromeTracer struct {
	instance *godet.RemoteDebugger
}

func NewChromeTracer(chrome *godet.RemoteDebugger) *chromeTracer {
	return &chromeTracer{
		instance: chrome,
	}
}

func (ct *chromeTracer) GetTrace(url *url.URL) ([]*redirect, error) {
	var frameId string
	var redirects []*redirect

	rawRedirects := make(map[string][]godet.Params)

	err := ct.instance.EnableRequestInterception(true)
	if err != nil {
		return redirects, errors.New(fmt.Sprintf("EnableRequestInterception failed. %s", err))
	}

	// get list of open tabs
	//tabs, err := ct.instance.TabList("")
	//if err != nil {
	//	return redirects, errors.New(fmt.Sprintf("TabList failed. %s", err))
	//}

	ct.instance.CallbackEvent("Network.requestWillBeSent", func(params godet.Params) {
		if _, ok := params["redirectResponse"]; ok && params["type"] == "Document" {
			rawRedirects[params["frameId"].(string)] = append(rawRedirects[params["frameId"].(string)], params)
		}
	})

	// create new tab
	tab, _ := ct.instance.NewTab("https://www.google.com")
	defer func() {
		err = ct.instance.CloseTab(tab)
		if err != nil {
			log.Error(errors.New(fmt.Sprintf("CloseTab failed. %s", err)))
		}
	}()

	// enable event processing
	//err = ct.instance.RuntimeEvents(true)
	//if err != nil {
	//	return redirects, errors.New(fmt.Sprintf("RuntimeEvents failed. %s", err))
	//}
	err = ct.instance.NetworkEvents(true)
	if err != nil {
		return redirects, errors.New(fmt.Sprintf("NetworkEvents failed. %s", err))
	}
	//err = ct.instance.PageEvents(true)
	//if err != nil {
	//	return redirects, errors.New(fmt.Sprintf("PageEvents failed. %s", err))
	//}
	//err = ct.instance.DOMEvents(true)
	//if err != nil {
	//	return redirects, errors.New(fmt.Sprintf("DOMEvents failed. %s", err))
	//}
	//err = ct.instance.LogEvents(true)
	//if err != nil {
	//	return redirects, errors.New(fmt.Sprintf("LogEvents failed. %s", err))
	//}

	// navigate in existing tab
	err = ct.instance.ActivateTab(tab)
	if err != nil {
		return redirects, errors.New(fmt.Sprintf("ActivateTab failed. %s", err))
	}

	// re-enable events when changing active tab
	err = ct.instance.AllEvents(true) // enable all events
	if err != nil {
		return redirects, errors.New(fmt.Sprintf("AllEvents failed. %s", err))
	}

	frameId, err = ct.instance.Navigate(url.String())
	if err != nil {
		return redirects, errors.New(fmt.Sprintf("Navigate failed. %s", err))
	}

	time.Sleep(time.Duration(time.Second * 5))

	if len(rawRedirects) == 0 {
		return redirects, errors.New("No redirects found")
	}

	if frameId == "" {
		return redirects, errors.New("Invalid mainframe id")
	}

	if rawRedirects, ok := rawRedirects[frameId]; ok {
		for _, rawRedirect := range rawRedirects {
			redirect, err := parseRedirectFromRaw(rawRedirect)
			if err != nil {
				return redirects, errors.New(fmt.Sprintf("An error during parsing redirects. %s", err))
			}

			redirects = append(redirects, redirect)
		}
	} else {
		return redirects, errors.New("No redirects found for mainframe")
	}

	return redirects, nil
}

func (ct *chromeTracer) Screenshot(url *url.URL, size *screenSize, path string) error {
	err := ct.instance.EnableRequestInterception(true)
	if err != nil {
		return errors.New(fmt.Sprintf("EnableRequestInterception failed. %s", err))
	}

	// create new tab
	tab, _ := ct.instance.NewTab(url.String())
	defer func() {
		err = ct.instance.CloseTab(tab)
		if err != nil {
			log.Error(errors.New(fmt.Sprintf("CloseTab failed. %s", err)))
		}
	}()

	// navigate in existing tab
	err = ct.instance.ActivateTab(tab)
	if err != nil {
		return errors.New(fmt.Sprintf("ActivateTab failed. %s", err))
	}

	_, err = ct.instance.Navigate(url.String())
	if err != nil {
		return errors.New(fmt.Sprintf("Navigate failed. %s", err))
	}

	time.Sleep(time.Duration(time.Second * 5))

	err = ct.instance.SetDeviceMetricsOverride(size.Width, size.Height, 0, false, false)
	if err != nil {
		return errors.New(fmt.Sprintf("Set screen size error: %s", err))
	}

	err = ct.instance.SetVisibleSize(size.Width, size.Height)
	if err != nil {
		return errors.New(fmt.Sprintf("Set visibility size error: %s", err))
	}

	//TODO: full page screenshot

	// take a screenshot
	err = ct.instance.SaveScreenshot(path, 0644, 100, true)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot capture screenshot: %s", err))
	}
	//time.Sleep(time.Second)

	return nil
}

func parseRedirectFromRaw(rawRedirect godet.Params) (*redirect, error) {
	if _, ok := rawRedirect["redirectResponse"]; !ok {
		return nil, errors.New("Invalid redirect. redirectResponse param not exists")
	}
	if _, ok := rawRedirect["request"]; !ok {
		return nil, errors.New("Invalid redirect. request param not exists")
	}

	redirectResponse := rawRedirect.Map("redirectResponse")
	request := rawRedirect.Map("request")

	to, err := url.Parse(rawRedirect.String("documentURL"))
	if err != nil {
		return nil, errors.New("Invalid redirect To url")
	}

	if _, ok := redirectResponse["url"]; !ok {
		return nil, errors.New("Invalid redirect. redirectResponse param url not exists")
	}
	fromUrl := redirectResponse["url"].(string)
	from, err := url.Parse(fromUrl)
	if err != nil {
		return nil, errors.New("Invalid redirect From url")
	}

	var cookies []*http.Cookie
	responseHeaders := http.Header{}
	if _, ok := redirectResponse["headers"]; !ok {
		return nil, errors.New("Invalid redirect. redirectResponse param headers not exists")
	}
	headers := redirectResponse["headers"].(map[string]interface{})
	for index, header := range headers {
		responseHeaders.Add(index, header.(string))
		if index == "Set-Cookie" {
			cookies = parseCookies(header.(string))
		}
	}

	requestHeaders := http.Header{}
	if _, ok := request["headers"]; !ok {
		return nil, errors.New("Invalid redirect. request param headers not exists")
	}

	headers = request["headers"].(map[string]interface{})
	for index, header := range headers {
		requestHeaders.Add(index, header.(string))
	}

	status := int(redirectResponse["status"].(float64))

	initiator := rawRedirect.Map("initiator")["type"].(string)

	redirect := NewRedirect(from, to, &requestHeaders, &responseHeaders, cookies, status, initiator)

	return redirect, nil
}

func parseCookies(s string) []*http.Cookie {
	return (&http.Response{Header: http.Header{"Set-Cookie": {s}}}).Cookies()
}

