package tracer

import "net/url"

type ITracer interface {
	GetTrace(url *url.URL) ([]*redirect, error)
	Screenshot(url *url.URL, size *screenSize, path string) error
}
