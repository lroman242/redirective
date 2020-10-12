package presenter

import "github.com/lroman242/redirective/domain"

// TracePresenter interface represent presenter service related for trace results
type TracePresenter interface {
	ResponseTraceResults(*domain.TraceResults) *domain.TraceResults
	ResponseScreenshot(string) string
}
