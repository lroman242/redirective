// Package domain describe base models used in application
package domain

import (
	"net/http"
	"net/url"
	"time"
)

// Redirect type represent http redirect.
type Redirect struct {
	From               *url.URL               `json:"from"`
	To                 *url.URL               `json:"to"`
	RequestHeaders     *http.Header           `json:"request_headers"`
	ResponseHeaders    *http.Header           `json:"response_headers"`
	Cookies            []*http.Cookie         `json:"cookies"`
	Status             int                    `json:"status"`
	Initiator          string                 `json:"initiator"`
	OtherInfo          map[string]interface{} `json:"other_info"`
	ScreenshotFileName string                 `json:"screenshot,omitempty"`
}

// NewRedirect combine data from http request and response to create
// new `Redirect` instance.
func NewRedirect(from, to *url.URL, requestHeaders, responseHeaders *http.Header, cookies []*http.Cookie, status int, initiator string) *Redirect {
	return &Redirect{
		From:            from,
		To:              to,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
		Cookies:         cookies,
		Status:          status,
		Initiator:       initiator,
	}
}

// JSONRedirect used to transform Redirect type into json string.
type JSONRedirect struct {
	From               string                 `json:"from"`
	To                 string                 `json:"to"`
	RequestHeaders     map[string]string      `json:"request_headers"`
	ResponseHeaders    map[string]string      `json:"response_headers"`
	Cookies            []*JSONCookie          `json:"cookies"`
	Status             int                    `json:"status"`
	Initiator          string                 `json:"initiator"`
	OtherInfo          map[string]interface{} `json:"other_info"`
	ScreenshotFileName string                 `json:"screenshot,omitempty"`
}

// NewJSONRedirects transform slice of `Redirect`s to slice of `jsonRedirect`s
// using NewJSONRedirect function.
func NewJSONRedirects(redirects []*Redirect) []*JSONRedirect {
	jsonRedirects := make([]*JSONRedirect, 0, len(redirects))

	for _, r := range redirects {
		jsonRedirects = append(jsonRedirects, NewJSONRedirect(r))
	}

	return jsonRedirects
}

// NewJSONRedirect function process `Redirect` to create
// `jsonRedirect` instance which can be marshaled to json.
func NewJSONRedirect(r *Redirect) *JSONRedirect {
	rRequestHeaders := make(map[string]string)
	for k := range *r.RequestHeaders {
		rRequestHeaders[k] = r.RequestHeaders.Get(k)
	}

	rResponseHeaders := make(map[string]string)
	for k := range *r.ResponseHeaders {
		rResponseHeaders[k] = r.ResponseHeaders.Get(k)
	}

	return &JSONRedirect{
		From:               r.From.String(),
		To:                 r.To.String(),
		RequestHeaders:     rRequestHeaders,
		ResponseHeaders:    rResponseHeaders,
		Cookies:            NewJSONCookies(r.Cookies),
		Status:             r.Status,
		Initiator:          r.Initiator,
		OtherInfo:          r.OtherInfo,
		ScreenshotFileName: r.ScreenshotFileName,
	}
}

// JSONCookie transform http.Cookie into json string.
type JSONCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string    `json:"path"`        // optional
	Domain     string    `json:"domain"`      // optional
	Expires    time.Time `json:"expires"`     // optional
	RawExpires string    `json:"raw_expires"` // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxA	ge>0 means Max-Age attribute present and given in seconds
	MaxAge   int      `json:"max_age"`
	Secure   bool     `json:"secure"`
	HTTPOnly bool     `json:"http_only"`
	Raw      string   `json:"raw"`
	Unparsed []string `json:"unparsed"` // Raw text of unparsed attribute-value pairs
}

// NewJSONCookies convert slice of `http.Cookies` to the slice of `jsonCookies` type
// using NewJSONCookie function.
func NewJSONCookies(cookies []*http.Cookie) []*JSONCookie {
	jsonCookies := make([]*JSONCookie, 0, len(cookies))

	for _, c := range cookies {
		jsonCookies = append(jsonCookies, NewJSONCookie(c))
	}

	return jsonCookies
}

// NewJSONCookie function is used to transform `http.Cookie` instance
// to a `JSONCookie`, which contains custom json marshal rules.
func NewJSONCookie(cookie *http.Cookie) *JSONCookie {
	return &JSONCookie{
		Name:       cookie.Name,
		Value:      cookie.Value,
		Path:       cookie.Path,
		Domain:     cookie.Domain,
		Expires:    cookie.Expires,
		RawExpires: cookie.RawExpires,
		MaxAge:     cookie.MaxAge,
		Secure:     cookie.Secure,
		HTTPOnly:   cookie.HttpOnly,
		Raw:        cookie.Raw,
		Unparsed:   cookie.Unparsed,
	}
}
