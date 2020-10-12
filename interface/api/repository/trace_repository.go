package repository

import (
	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/infrastructure/storage"
)

// TraceRepository interface represent repository to work with trace results
type TraceRepository interface {
	SaveTraceResults(*domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}

type traceRepository struct {
	storage storage.Storage
}

func NewTraceRepository(storage storage.Storage) TraceRepository {
	return &traceRepository{storage}
}

func (tr *traceRepository) SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error) {
	return tr.storage.SaveTraceResults(traceResults)
}

func (tr *traceRepository) FindTraceResults(id interface{}) (*domain.TraceResults, error) {
	return tr.storage.FindTraceResults(id)
}
