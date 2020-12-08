package presenter_test

import (
	"testing"

	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/interface/api/presenter"
)

const (
	testDomain   = `redirective.net`
	testProtocol = `http`
	testFileName = `screenshot.png`
)

func TestTracePresenter_ResponseTraceResults(t *testing.T) {
	tr := &domain.TraceResults{
		ID:         "someResultsId",
		Redirects:  make([]*domain.Redirect, 0),
		Screenshot: testFileName,
		URL:        "",
	}

	tp := presenter.NewTracePresenter(testDomain, testProtocol)
	result := tp.ResponseTraceResults(tr)

	if result.URL != testProtocol+"://"+testDomain+"/api/find/"+result.ID.(string) {
		t.Error("invalid result url received")
	}

	if result.Screenshot != testProtocol+"://"+testDomain+"/screenshots/"+testFileName {
		t.Error("invalid result screenshot url received")
	}
}

func TestTracePresenter_ResponseScreenshot(t *testing.T) {
	tp := presenter.NewTracePresenter(testDomain, testProtocol)
	result := tp.ResponseScreenshot(testFileName)

	if result != testProtocol+"://"+testDomain+"/screenshots/"+testFileName {
		t.Error("invalid result received")
	}
}
