package interactor

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/mocks"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestTraceInteractor_Trace_Success(t *testing.T) {
	assetsFolderPath := "assets/screenshots/"
	u, err := url.Parse("http://ssyoutube.com")
	if err != nil {
		t.Errorf("Cannot parse error: %s\n", err)
	}

	expectedId := 13
	tr := &domain.TraceResults{}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	tracer.EXPECT().Trace(u, gomock.Any()).Times(1).Return(tr, nil)

	traceRepository := mocks.NewMockTraceRepository(mockCtrl)
	traceRepository.EXPECT().SaveTraceResults(tr).Times(1).Return(expectedId, nil)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)
	//check if TracePresenter called correctly
	tracePresenter.EXPECT().ResponseTraceResults(tr).Times(1).DoAndReturn(func(tr *domain.TraceResults) interface{} {
		if tr.ID != expectedId {
			t.Error("invalid results ID provided to TracePresenter")
		}

		tr.ID = expectedId + 5

		return tr
	})

	logger := mocks.NewMockLogger(mockCtrl)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)

	results, err := ti.Trace(u, assetsFolderPath)
	if err != nil {
		t.Errorf("Invalid response received. trace failed. Error: %s\n", err)
	}

	if results.ID != (expectedId + 5) {
		t.Error("unexpected result received")
	}
}

func TestTraceInteractor_Trace_TracerError(t *testing.T) {
	assetsFolderPath := "assets/screenshots/"
	u, err := url.Parse("http://ssyoutube.com")
	if err != nil {
		t.Errorf("Cannot parse error: %s\n", err)
	}

	expectedError := errors.New("expected Tracer error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	tracer.EXPECT().Trace(u, gomock.Any()).Times(1).Return(nil, expectedError)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)

	traceRepository := mocks.NewMockTraceRepository(mockCtrl)

	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError).Times(1)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)

	results, err := ti.Trace(u, assetsFolderPath)
	if results != nil {
		t.Error("unexpected response received. no results expected")
	}
	if err == nil {
		t.Error("unexpected response received. an error expected")
	} else {
		if !reflect.DeepEqual(err, expectedError) {
			t.Errorf("unexpected error received: %s\n", err)
		}
	}
}

func TestTraceInteractor_Trace_TraceRepository_SaveTraceResults_Error(t *testing.T) {
	assetsFolderPath := "assets/screenshots/"
	u, err := url.Parse("http://ssyoutube.com")
	if err != nil {
		t.Errorf("Cannot parse error: %s\n", err)
	}

	expectedError := errors.New("expected TracerRepository error")

	tr := &domain.TraceResults{}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	tracer.EXPECT().Trace(u, gomock.Any()).Times(1).Return(tr, nil)

	traceRepository := mocks.NewMockTraceRepository(mockCtrl)
	traceRepository.EXPECT().SaveTraceResults(tr).Times(1).Return(nil, expectedError)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)

	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError).Times(1)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)

	results, err := ti.Trace(u, assetsFolderPath)
	if results == nil {
		t.Error("unexpected response received. some results expected")
	}
	if err == nil {
		t.Error("unexpected response received. an error expected")
	} else {
		if !reflect.DeepEqual(err, expectedError) {
			t.Errorf("unexpected error received: %s\n", err)
		}
	}
}

func TestTraceInteractor_FindTrace_Success(t *testing.T) {
	expectedResultsID := 17
	expectedScreenshotPath := "assets/screenshots/" + randomScreenshotFileName("png")

	tr := &domain.TraceResults{
		ID: expectedResultsID,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	traceRepository := mocks.NewMockTraceRepository(mockCtrl)
	traceRepository.EXPECT().FindTraceResults(expectedResultsID).Times(1).Return(tr, nil)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)
	tracePresenter.EXPECT().ResponseTraceResults(tr).Times(1).DoAndReturn(func(tr *domain.TraceResults) interface{} {
		if tr.ID != expectedResultsID {
			t.Error("invalid results ID provided to TracePresenter")
		}
		tr.Screenshot = expectedScreenshotPath

		return tr
	})

	logger := mocks.NewMockLogger(mockCtrl)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)
	results, err := ti.FindTrace(expectedResultsID)
	if err != nil {
		t.Errorf("unexpected error received: %s\n", err)
	}
	if results.ID != expectedResultsID {
		t.Error("unexpected results received. wrong results.ID")
	}
	if results.Screenshot != expectedScreenshotPath {
		t.Error("unexpected results received. wrong results.Screenshot")
	}
}

func TestTraceInteractor_FindTrace_Repository_FindTraceResults_Error(t *testing.T) {
	expectedResultsID := 21
	expectedError := errors.New("expected TracerRepository error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	traceRepository := mocks.NewMockTraceRepository(mockCtrl)
	traceRepository.EXPECT().FindTraceResults(expectedResultsID).Times(1).Return(nil, expectedError)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)

	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError).Times(1)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)
	results, err := ti.FindTrace(expectedResultsID)
	if results != nil {
		t.Error("unexpected response received. no results expected")
	}
	if err == nil {
		t.Error("unexpected response received. an error expected")
	} else {
		if !reflect.DeepEqual(err, expectedError) {
			t.Errorf("unexpected error received: %s\n", err)
		}
	}
}

func TestTraceInteractor_Screenshot_Success(t *testing.T) {
	assetsFolderPath := "assets/screenshots/"
	screenshotPath := assetsFolderPath + "someRandomFileName.png"
	u, err := url.Parse("http://ssyoutube.com")
	if err != nil {
		t.Errorf("Cannot parse error: %s\n", err)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	tracer.EXPECT().Screenshot(u, gomock.Any(), gomock.Any()).Times(1).Return(nil)

	traceRepository := mocks.NewMockTraceRepository(mockCtrl)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)
	//check if TracePresenter called correctly
	tracePresenter.EXPECT().ResponseScreenshot(gomock.Any()).Times(1).DoAndReturn(func(screenshot string) interface{} {
		screenshot = screenshotPath
		return strings.Replace(screenshot, assetsFolderPath, "https://redirective.net/screenshots/", 1)
	})

	logger := mocks.NewMockLogger(mockCtrl)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)

	path, err := ti.Screenshot(u, 1920, 1080, assetsFolderPath)
	if err != nil {
		t.Errorf("Invalid response received. trace failed. Error: %s\n", err)
	}

	if path != "https://redirective.net/screenshots/someRandomFileName.png" {
		t.Errorf("unexpected result received. %s", path)
	}
}

func TestTraceInteractor_Screenshot_Tracer_Screenshot_Error(t *testing.T) {
	assetsFolderPath := "assets/screenshots/"
	expectedError := errors.New("expected Tracer error")

	u, err := url.Parse("http://ssyoutube.com")
	if err != nil {
		t.Errorf("Cannot parse error: %s\n", err)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracer := mocks.NewMockTracer(mockCtrl)
	tracer.EXPECT().Screenshot(u, gomock.Any(), gomock.Any()).Times(1).Return(expectedError)

	traceRepository := mocks.NewMockTraceRepository(mockCtrl)

	tracePresenter := mocks.NewMockTracePresenter(mockCtrl)
	logger := mocks.NewMockLogger(mockCtrl)
	logger.EXPECT().Error(expectedError)

	ti := NewTraceInteractor(tracer, tracePresenter, traceRepository, logger)

	path, err := ti.Screenshot(u, 1920, 1080, assetsFolderPath)
	if err == nil {
		t.Error("unexpected response received. an error expected")
	} else {
		if !reflect.DeepEqual(err, expectedError) {
			t.Errorf("unexpected error received: %s\n", err)
		}
	}

	if path != "" {
		t.Errorf("unexpected result received. %s", path)
	}
}
