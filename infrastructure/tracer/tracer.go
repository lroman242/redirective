// Package tracer implements types and methods to trace http requests
package tracer

import (
	"github.com/lroman242/redirective/domain"
	"net/url"
)

// Tracer interface represent required list of function for http tracers
type Tracer interface {
	Trace(url *url.URL, fileName string) ([]*domain.Redirect, error)
	Screenshot(url *url.URL, size *domain.ScreenSize, path string) error
}
