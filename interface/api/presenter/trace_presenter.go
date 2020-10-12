package presenter

import (
	"github.com/lroman242/redirective/domain"
)

// TracePresenter interface represent presenter service related for trace results
type TracePresenter interface {
	ResponseTraceResults(*domain.TraceResults) *domain.TraceResults
	ResponseScreenshot(string) string
}

type tracePresenter struct {
}

// NewTracePresenter will construct TracePresenter implementation
func NewTracePresenter() TracePresenter {
	return &tracePresenter{}
}

// ResponseTraceResults will provide updated trace results
func (t *tracePresenter) ResponseTraceResults(results *domain.TraceResults) *domain.TraceResults {
	// TODO: update trace result screenshot path to url

	return results
}

// ResponseScreenshot will change screenshot path to url
func (t *tracePresenter) ResponseScreenshot(filename string) string {
	return "https://redirective.net/assets/" + filename
}
