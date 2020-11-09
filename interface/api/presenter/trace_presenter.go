package presenter

import (
	"github.com/lroman242/redirective/domain"
)

// TracePresenter interface represent presenter service related for trace results
type TracePresenter interface {
	ResponseTraceResults(*domain.TraceResults) *domain.TraceResults
	ResponseScreenshot(string) string
}

// tracePresenter is implementation of TracePresenter interface
type tracePresenter struct {
	appDomain string
	protocol  string
}

// NewTracePresenter will construct TracePresenter implementation
func NewTracePresenter(appDomain string, protocol string) TracePresenter {
	return &tracePresenter{
		appDomain: appDomain,
		protocol:  protocol,
	}
}

// ResponseTraceResults will provide updated trace results
func (t *tracePresenter) ResponseTraceResults(results *domain.TraceResults) *domain.TraceResults {
	results.Screenshot = t.ResponseScreenshot(results.Screenshot)
	results.URL = t.protocol + "://" + t.appDomain + "/api/find/" + results.ID.(string)

	return results
}

// ResponseScreenshot will change screenshot path to url
func (t *tracePresenter) ResponseScreenshot(filename string) string {
	return t.protocol + "://" + t.appDomain + "/screenshots/" + filename
}
