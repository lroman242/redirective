// Package tracer implements types and methods to trace http requests
package tracer

import "net/url"

// Tracer interface represent required list of function for http tracers
type Tracer interface {
	GetTrace(url *url.URL) ([]*Redirect, error)
	Screenshot(url *url.URL, size *ScreenSize, path string) error
}
