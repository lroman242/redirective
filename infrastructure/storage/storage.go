package storage

import "github.com/lroman242/redirective/domain"

//go:generate mockgen -package=mocks -destination=mocks/mock_storage.go -source=infrastructure/storage/storage.go Storage

// Storage interface describe storage instance used to store some data
type Storage interface {
	SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}
