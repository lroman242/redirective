package repository

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/mocks"
	"reflect"
	"testing"
)

func TestTraceRepository_FindTraceResults(t *testing.T) {
	testResultsID := "SomeID"

	testResults := &domain.TraceResults{
		ID: testResultsID,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().FindTraceResults(testResultsID).Times(1).Return(testResults, nil)

	tr := NewTraceRepository(storage)
	result, err := tr.FindTraceResults(testResultsID)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(result, testResults) {
		t.Error("invalid results received")
	}
}

func TestTraceRepository_FindTraceResults_Error(t *testing.T) {
	expectedError := errors.New("expected error")
	testResultsID := "SomeID"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().FindTraceResults(testResultsID).Times(1).Return(nil, expectedError)

	tr := NewTraceRepository(storage)
	result, err := tr.FindTraceResults(testResultsID)
	if err == nil {
		t.Error("an error expected")
	}
	if result != nil {
		t.Error("no results expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Error("wrong error received")
	}
}

func TestTraceRepository_SaveTraceResults(t *testing.T) {
	testResultsID := "SomeID"

	testResults := &domain.TraceResults{}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().SaveTraceResults(testResults).Times(1).Return(testResultsID, nil)

	tr := NewTraceRepository(storage)
	ID, err := tr.SaveTraceResults(testResults)
	if err != nil {
		t.Error(err)
	}
	if ID.(string) != testResultsID {
		t.Error("invalid results ID received")
	}
}

func TestTraceRepository_SaveTraceResults_Error(t *testing.T) {
	expectedError := errors.New("expected error")

	testResults := &domain.TraceResults{}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storage := mocks.NewMockStorage(mockCtrl)
	storage.EXPECT().SaveTraceResults(testResults).Times(1).Return(nil, expectedError)

	tr := NewTraceRepository(storage)
	ID, err := tr.SaveTraceResults(testResults)
	if err == nil {
		t.Error("an error expected")
	}
	if ID != nil {
		t.Error("no id expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Error("wrong error received")
	}
}
