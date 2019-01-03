package tracer

import (
	"net/http"
	"net/url"
	"time"
)

type redirect struct {
	From            *url.URL               `json:"from"`
	To              *url.URL               `json:"to"`
	RequestHeaders  *http.Header           `json:"request_headers"`
	ResponseHeaders *http.Header           `json:"response_headers"`
	Cookies         []*http.Cookie         `json:"cookies"`
	Status          int                    `json:"status"`
	Initiator       string                 `json:"initiator"`
	OtherInfo       map[string]interface{} `json:"other_info"`
}

func NewRedirect(from, to *url.URL, requestHeaders, responseHeaders *http.Header, cookies []*http.Cookie, status int, initiator string) *redirect {
	return &redirect{
		From:            from,
		To:              to,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
		Cookies:         cookies,
		Status:          status,
		Initiator:       initiator,
	}
}

type jsonRedirect struct {
	From            string                 `json:"from"`
	To              string                 `json:"to"`
	RequestHeaders  map[string]string      `json:"request_headers"`
	ResponseHeaders map[string]string      `json:"response_headers"`
	Cookies         []*jsonCookie          `json:"cookies"`
	Status          int                    `json:"status"`
	Initiator       string                 `json:"initiator"`
	OtherInfo       map[string]interface{} `json:"other_info"`
}

func NewJSONRedirects(redirects []*redirect) []*jsonRedirect {
	var jsonRedirects []*jsonRedirect

	for _, r := range redirects {
		jsonRedirects = append(jsonRedirects, NewJSONRedirect(r))
	}

	return jsonRedirects
}

func NewJSONRedirect(r *redirect) *jsonRedirect {
	rRequestHeaders := make(map[string]string)
	rResponseHeaders := make(map[string]string)

	for k := range *r.RequestHeaders {
		rRequestHeaders[k] = r.RequestHeaders.Get(k)
	}
	for k := range *r.ResponseHeaders {
		rResponseHeaders[k] = r.RequestHeaders.Get(k)
	}

	return &jsonRedirect{
		From:            r.From.String(),
		To:              r.To.String(),
		RequestHeaders:  rRequestHeaders,
		ResponseHeaders: rResponseHeaders,
		Cookies:         NewJSONCookies(r.Cookies),
		Status:          r.Status,
		Initiator:       r.Initiator,
		OtherInfo:       r.OtherInfo,
	}
}

type jsonCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string    `json:"path"`        // optional
	Domain     string    `json:"domain"`      // optional
	Expires    time.Time `json:"expires"`     // optional
	RawExpires string    `json:"raw_expires"` // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int      `json:"max_age"`
	Secure   bool     `json:"secure"`
	HttpOnly bool     `json:"http_only"`
	Raw      string   `json:"raw"`
	Unparsed []string `json:"unparsed"` // Raw text of unparsed attribute-value pairs
}

func NewJSONCookies(cookies []*http.Cookie) []*jsonCookie {
	var jsonCookies []*jsonCookie

	for _, c := range cookies {
		jsonCookies = append(jsonCookies, NewJSONCookie(c))
	}

	return jsonCookies
}

func NewJSONCookie(cookie *http.Cookie) *jsonCookie {
	return &jsonCookie{
		Name:       cookie.Name,
		Value:      cookie.Value,
		Path:       cookie.Path,
		Domain:     cookie.Domain,
		Expires:    cookie.Expires,
		RawExpires: cookie.RawExpires,
		MaxAge:     cookie.MaxAge,
		Secure:     cookie.Secure,
		HttpOnly:   cookie.HttpOnly,
		Raw:        cookie.Raw,
		Unparsed:   cookie.Unparsed,
	}
}
