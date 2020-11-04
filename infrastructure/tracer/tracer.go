// Package tracer implements types and methods to trace http requests
package tracer

import (
	"github.com/lroman242/redirective/domain"
	"net/url"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_tracer.go -source=infrastructure/tracer/tracer.go Tracer

// Tracer interface represent required list of function for http tracers
type Tracer interface {
	Trace(url *url.URL, path string) (*domain.TraceResults, error)
	Screenshot(url *url.URL, size *ScreenSize, path string) error
}
