package tracer

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lroman242/redirective/domain"
	"github.com/raff/godet"
)

const (
	setCookieHeaderName = "set-cookie"
	documentParamName   = "Document"
)

const (
	screenshotDelay = 3
)

// ChromeRemoteDebuggerInterface implements API to work with browser debugger.
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
	SetUserAgent(userAgent string) error
	Close() error
}

// ChromeTracer interface extend Tracer interface
// and describe additional functions required for chrome browser.
type ChromeTracer interface {
	Tracer
	ChromeProcess() *os.Process
	Close() error
}

// chromeTracer represent tracer based on google chrome debugging tools.
type chromeTracer struct {
	chromePort int
	// chromePath             string
	size                   *ScreenSize
	screenshotsStoragePath string
	chromeProcess          *os.Process
}

func (ct *chromeTracer) initChromeRemoteDebugger() (ChromeRemoteDebuggerInterface, error) {
	remote, err := godet.Connect("localhost:"+strconv.Itoa(ct.chromePort), false)
	if err != nil {
		log.Fatalf("Cannot connect to remote debugger: %s\n", err)

		return nil, err
	}

	remote.CallbackEvent("RemoteDebugger.disconnected", func(params godet.Params) {
		log.Println("Remote disconnected")
		panic("Remote disconnected")
	})

	return remote, nil
}

// Close method will stop google-chrome process.
func (ct *chromeTracer) Close() error {
	//if err := ct.chromeProcess.Kill(); err != nil {
	//	log.Fatalf("Close error: %s\n", err)
	//
	//	return err
	//}

	return nil
}

// ChromeProcess get google-chrome process.
func (ct *chromeTracer) ChromeProcess() *os.Process {
	return ct.chromeProcess
}

// NewChromeTracer create new chrome tracer instance.
func NewChromeTracer(size *ScreenSize, screenshotsStoragePath string) ChromeTracer {
	// /usr/bin/google-chrome --addr=localhost --port=9222 --remote-debugging-port=9222 --remote-debugging-address=0.0.0.0 --disable-extensions --disable-gpu --headless --hide-scrollbars --no-first-run --no-sandbox --user-agent="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/77.0.3854.3 Chrome/77.0.3854.3 Safari/537.36"
	// TODO: get chrome port and path from func args
	//cmd := exec.Command("/usr/bin/google-chrome",
	//	"--addr=localhost",
	//	"--port=9222",
	//	"--remote-debugging-port=9222",
	//	"--remote-debugging-address=0.0.0.0",
	//	"--disable-extensions",
	//	"--disable-gpu",
	//	"--headless",
	//	"--hide-scrollbars",
	//	"--no-first-run",
	//	"--no-sandbox",
	//	"--user-agent=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	//
	//cmd.Stdout = os.Stdout
	//
	//err := cmd.Start()
	//if err != nil {
	//	panic(err)
	//}

	ct := &chromeTracer{
		chromePort: 9222,
		// chromePath:             path,
		size:                   size,
		screenshotsStoragePath: screenshotsStoragePath,
		// chromeProcess:          cmd.Process,
	}

	return ct
}

func (ct *chromeTracer) traceURL(debugger ChromeRemoteDebuggerInterface, url *url.URL, redirects, responses *map[string][]godet.Params, filePath string) (string, error) {
	frameID := ""

	err := debugger.EnableRequestInterception(true)
	if err != nil {
		return frameID, &EnableRequestInterceptionError{err: err}
	}

	debugger.CallbackEvent("Network.requestWillBeSent", func(params godet.Params) {
		if _, ok := params["redirectResponse"]; ok && params["type"] == documentParamName {
			(*redirects)[params["frameId"].(string)] = append((*redirects)[params["frameId"].(string)], params)
		}
	})
	debugger.CallbackEvent("Network.responseReceived", func(params godet.Params) {
		if params["type"] == documentParamName {
			(*responses)[params["frameId"].(string)] = append((*responses)[params["frameId"].(string)], params)
		}
	})

	frameID, err = newTab(debugger, url, ct.size, filePath)
	if err != nil {
		return frameID, &NewTabError{err: err}
	}

	return frameID, nil
}

// Trace parse redirect trace path for provided url.
func (ct *chromeTracer) Trace(url *url.URL, filePath string) (*domain.TraceResults, error) {
	debugger, err := ct.initChromeRemoteDebugger()
	if err != nil {
		log.Panicf("Remote debugger is not awailable: %s", err)
	}

	defer func() {
		err = debugger.Close()
		if err != nil {
			log.Printf("Cannot close debugger: %s \n\n", err)
		}
	}()

	var redirects []*domain.Redirect

	rawRedirects := make(map[string][]godet.Params)
	rawResponses := make(map[string][]godet.Params)

	frameID, err := ct.traceURL(debugger, url, &rawRedirects, &rawResponses, filePath)
	if err != nil {
		return nil, err
	}

	if frameID == "" {
		return nil, &InvalidFrameIDError{}
	}

	//if len(rawRedirects) == 0 {
	//	return nil, errors.New(`no data parsed`)
	//}

	if rawRedirects, ok := rawRedirects[frameID]; ok {
		for _, rawRedirect := range rawRedirects {
			redirect, err := parseRedirectFromRaw(rawRedirect)
			if err != nil {
				return nil, &RedirectParseError{err: err}
			}

			redirects = append(redirects, redirect)
		}
	} /* else {
		return redirects, errors.New("No redirects found for mainframe")
	}*/

	if rawResponse, ok := rawResponses[frameID]; ok {
		response, err := parseMainResponseFromRaw(rawResponse[len(rawResponse)-1])
		if err != nil {
			return nil, &ResponseParseError{err: err}
		}

		redirects = append(redirects, response)
	} else {
		return nil, &NoResponseError{}
	}

	return &domain.TraceResults{
		Redirects:  redirects,
		Screenshot: filePath,
	}, nil
}

// Screenshot function makes a final page screen capture.
func (ct *chromeTracer) Screenshot(url *url.URL, size *ScreenSize, filePath string) error {
	debugger, err := ct.initChromeRemoteDebugger()
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = debugger.Close()
	}()

	err = debugger.EnableRequestInterception(true)
	if err != nil {
		return &EnableRequestInterceptionError{err: err}
	}

	_, err = newTab(debugger, url, size, filePath)
	if err != nil {
		return &NewTabError{err: err}
	}

	return nil
}

func parseRedirectFromRaw(rawRedirect godet.Params) (*domain.Redirect, error) {
	if _, ok := rawRedirect["redirectResponse"]; !ok {
		return nil, &RedirectResponseNotExistsInRawDataError{}
	}

	if _, ok := rawRedirect["request"]; !ok {
		return nil, &RequestNotExistsInRawDataError{}
	}

	redirectResponse := rawRedirect.Map("redirectResponse")
	request := rawRedirect.Map("request")

	to, err := url.Parse(rawRedirect.String("documentURL"))
	if err != nil {
		return nil, &InvalidToURLDataError{}
	}

	if _, ok := redirectResponse["url"]; !ok {
		return nil, &URLParamNotExistsInRedirectDataError{}
	}

	from, err := url.Parse(redirectResponse["url"].(string))
	if err != nil {
		return nil, &InvalidFromURLDataError{}
	}

	if _, ok := redirectResponse["headers"]; !ok {
		return nil, &HeaderParamNotExistsInRedirectResponseDataError{}
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
		return nil, &InitiatorParamNotExistsInRedirectDataError{}
	}

	initiator := rawRedirect.Map("initiator")["type"].(string)

	return domain.NewRedirect(from, to, requestHeaders, &responseHeaders, cookies, status, initiator), nil
}

func parseHeadersFromRaw(request map[string]interface{}) (*http.Header, error) {
	requestHeaders := &http.Header{}

	if _, ok := request["headers"]; !ok {
		return requestHeaders, &HeaderParamNotExistsInRedirectDataError{}
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

func parseMainResponseFromRaw(rawResponses godet.Params) (*domain.Redirect, error) {
	if _, ok := rawResponses["response"]; !ok {
		return nil, &ResponseParamNotExistsInRedirectDataError{}
	}

	response := rawResponses.Map("response")

	if _, ok := response["url"]; !ok {
		return nil, &URLParamNotExistsInRedirectDataError{}
	}

	to, err := url.Parse(response["url"].(string))
	if err != nil {
		return nil, &InvalidToURLDataError{}
	}

	if _, ok := response["headers"]; !ok {
		return nil, &HeaderParamNotExistsInRedirectResponseDataError{}
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
		return nil, &HeaderParamNotExistsInRedirectDataError{}
	}

	requestHeaders := http.Header{}

	for index, header := range response["requestHeaders"].(map[string]interface{}) {
		requestHeaders.Add(index, header.(string))
	}

	status := int(response["status"].(float64))

	redirect := domain.NewRedirect(&url.URL{}, to, &requestHeaders, &responseHeaders, cookies, status, "")

	return redirect, nil
}

func newTab(debugger ChromeRemoteDebuggerInterface, url *url.URL, size *ScreenSize, filePath string) (string, error) {
	frameID := ``

	// create new tab
	tab, _ := debugger.NewTab(url.String())
	defer func(tab *godet.Tab) {
		err := debugger.CloseTab(tab)
		if err != nil {
			log.Println(&CloseTabError{err: err})
		}
	}(tab)

	err := debugger.NetworkEvents(true)
	if err != nil {
		return frameID, &NetworkEventsError{err: err}
	}

	// navigate in existing tab
	err = debugger.ActivateTab(tab)
	if err != nil {
		return frameID, &ActiveTabError{}
	}

	// re-enable events when changing active tab
	err = debugger.AllEvents(true) // enable all events
	if err != nil {
		return frameID, &AllEventsError{err: err}
	}

	err = debugger.SetDeviceMetricsOverride(size.Width, size.Height, 0, false, false)
	if err != nil {
		return frameID, &SetScreenSizeError{err: err}
	}

	err = debugger.SetVisibleSize(size.Width, size.Height)
	if err != nil {
		return frameID, &SetVisibilitySizeError{err: err}
	}

	frameID, err = debugger.Navigate(url.String())
	if err != nil {
		return frameID, &NavigateError{err: err}
	}

	time.Sleep(time.Second * screenshotDelay)

	// take a screenshot
	err = debugger.SaveScreenshot(filePath, 0o644, 100, true)
	if err != nil {
		return frameID, &CaptureScreenshotError{err: err}
	}

	return frameID, nil
}
