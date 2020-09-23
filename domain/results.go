package domain

// TraceResults describe composed data type that is used in response and storage
type TraceResults struct {
	Redirects  []*Redirect
	Screenshot string
}
