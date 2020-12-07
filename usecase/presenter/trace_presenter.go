// Package presenter contains interfaces to represent functions required to provide correct response format.
package presenter

import "github.com/lroman242/redirective/domain"

//go:generate mockgen -package=mocks -destination=mocks/mock_trace_presenter.go -source=usecase/presenter/trace_presenter.go TracePresenter

// TracePresenter interface represent presenter service related for trace results.
type TracePresenter interface {
	ResponseTraceResults(*domain.TraceResults) *domain.TraceResults
	ResponseScreenshot(string) string
}
