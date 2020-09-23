package storage

import "github.com/lroman242/redirective/domain"

type Storage interface {
	SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}
