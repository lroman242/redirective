package presenter

import (
	"github.com/lroman242/redirective/domain"
	"testing"
)

func TestTracePresenter_ResponseTraceResults(t *testing.T) {
	testDomain := "redirective.net"
	testProtocol := "http"
	testFileName := "screenshot.png"
	testResultsID := "someResultsId"

	tr := &domain.TraceResults{
		ID:         testResultsID,
		Redirects:  make([]*domain.Redirect, 0, 0),
		Screenshot: testFileName,
		URL:        "",
	}

	tp := NewTracePresenter(testDomain, testProtocol)
	result := tp.ResponseTraceResults(tr)

	if result.URL != testProtocol+"://"+testDomain+"/api/find/"+result.ID.(string) {
		t.Error("invalid result url received")
	}
	if result.Screenshot != testProtocol+"://"+testDomain+"/screenshots/"+testFileName {
		t.Error("invalid result screenshot url received")
	}
}

func TestTracePresenter_ResponseScreenshot(t *testing.T) {
	testDomain := "redirective.net"
	testProtocol := "http"
	testFileName := "screenshot.png"

	tp := NewTracePresenter(testDomain, testProtocol)
	result := tp.ResponseScreenshot(testFileName)

	if result != testProtocol+"://"+testDomain+"/screenshots/"+testFileName {
		t.Error("invalid result received")
	}
}
