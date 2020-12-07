package domain

// TraceResults describe composed data type that is used in response and storage.
type TraceResults struct {
	ID         interface{} `json:"id,omitempty"`
	Redirects  []*Redirect `json:"redirects"`
	Screenshot string      `json:"screenshot"`
	URL        string      `json:"url,omitempty"`
}
