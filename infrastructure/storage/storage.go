package storage

import "github.com/lroman242/redirective/domain"

// Storage interface describe storage instance used to store some data
type Storage interface {
	SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}
