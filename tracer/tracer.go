package tracer

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Tracer struct {
	Id                 uuid.UUID  `json:"id"`
	Url                url.URL    `json:"url"`
	Redirects          []Redirect `json:"redirects"`
	UserAgent          string     `json:"user_agent"`
	MaxRedirectsNumber int        `json:"max_redirects_number"`
}

type Redirect struct {
	StatusCode int
	Headers    map[string]interface{}
	Cookies    []*http.Cookie
	Body       []byte
	Source     *url.URL
	Target     url.URL
}

func NewTracer(url url.URL, userAgent string, maxRedirects int) *Tracer {
	trk := &Tracer{
		Id:                 uuid.NewV1(),
		Url:                url,
		MaxRedirectsNumber: maxRedirects,
		UserAgent:          userAgent,
	}

	return trk
}

func (t *Tracer) SetUrl(rawUrl string) error {
	if rawUrl == "" || !isUrl(rawUrl) {
		return errors.New("tracking url is not valid")
	}

	traceUrl, err := url.Parse(rawUrl)
	if err != nil {
		return errors.New("tracking url is not valid")
	}

	t.Url = *traceUrl

	return nil
}

func (t *Tracer) ProcessRedirects() error {
	nextURL := t.Url.String()

	var i int
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest("GET", nextURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", t.UserAgent)

	for i < t.MaxRedirectsNumber {
		resp, err := client.Do(req)
		if err != nil {
			return errors.New("sorry, an error occurred. please try again later")
		}

		redirect, err := parseResponse(resp)
		if err != nil {
			return err
		}
		t.Redirects = append(t.Redirects, redirect)

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			nextURL = resp.Header.Get("Location")
			i += 1
		} else {
			break
		}
	}

	return nil
}

func isUrl(trackUrl string) bool {
	_, err := url.ParseRequestURI(trackUrl)
	if err != nil {
		return false
	}
	return true
}

func parseResponse(resp *http.Response) (Redirect, error) {
	target, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return Redirect{}, errors.New("target url is not valid")
	}

	redirect := Redirect{
		Source:     resp.Request.URL,
		Target:     *target,
		StatusCode: resp.StatusCode,
		Cookies:    resp.Cookies(),
	}

	if len(resp.Header) > 0 {
		redirect.Headers = map[string]interface{}{}

		for k, h := range resp.Header {
			redirect.Headers[strings.ToLower(k)] = string(h[0])
		}
	}

	resp.Body.Read(redirect.Body)
	//_, err = resp.Body.Read(redirect.Body)
	//if err != nil {
	//	return Redirect{}, errors.New("Invalid response content")
	//}

	return redirect, nil
}
