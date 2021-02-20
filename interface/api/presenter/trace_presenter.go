// Package presenter contains classes to represent functions required to provide correct response format.
package presenter

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/lroman242/redirective/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TracePresenter interface represent presenter service related for trace results.
type TracePresenter interface {
	ResponseTraceResults(*domain.TraceResults) *domain.TraceResults
	ResponseScreenshot(string) string
}

// tracePresenter is implementation of TracePresenter interface.
type tracePresenter struct {
	appDomain string
	protocol  string
}

// NewTracePresenter will construct TracePresenter implementation.
func NewTracePresenter(appDomain string, protocol string) TracePresenter {
	return &tracePresenter{
		appDomain: appDomain,
		protocol:  protocol,
	}
}

// ResponseTraceResults will provide updated trace results.
func (t *tracePresenter) ResponseTraceResults(results *domain.TraceResults) *domain.TraceResults {
	results.Screenshot = t.ResponseScreenshot(results.Screenshot)
	results.URL = t.protocol + "://" + t.appDomain + "/api/find/" + t.getIDString(results.ID)

	return results
}

// ResponseScreenshot will change screenshot path to url.
func (t *tracePresenter) ResponseScreenshot(filename string) string {
	return t.protocol + "://" + t.appDomain + "/screenshots/" + filepath.Base(filename)
}

// getIDString function represent TraceResults.ID as string.
func (t *tracePresenter) getIDString(ID interface{}) string {
	switch tp := ID.(type) {
	case int:
		return strconv.Itoa(tp)
	case primitive.ObjectID:
		return tp.Hex()
	case *primitive.ObjectID:
		return tp.Hex()
	case string:
		return tp
	case fmt.Stringer:
		return tp.String()
	default:
		return fmt.Sprintf("%v", tp)
	}
}
