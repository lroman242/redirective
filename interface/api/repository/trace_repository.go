// Package repository contains classes that work with storage
package repository

import (
	"github.com/lroman242/redirective/domain"
	"github.com/lroman242/redirective/infrastructure/storage"
)

// TraceRepository interface represent repository to work with trace results.
type TraceRepository interface {
	SaveTraceResults(*domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}

// traceRepository is implementation of TraceRepository interface.
type traceRepository struct {
	storage storage.Storage
}

// NewTraceRepository function will build instance of TraceRepository.
func NewTraceRepository(storage storage.Storage) TraceRepository {
	return &traceRepository{storage}
}

// SaveTraceResults save TracerResults into storage.
func (tr *traceRepository) SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error) {
	return tr.storage.SaveTraceResults(traceResults)
}

// FindTraceResults will find and return domain.TracerResults from storage.
func (tr *traceRepository) FindTraceResults(id interface{}) (*domain.TraceResults, error) {
	return tr.storage.FindTraceResults(id)
}
