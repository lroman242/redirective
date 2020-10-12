package repository

import "github.com/lroman242/redirective/domain"

// TraceRepository interface represent repository to work with trace results
type TraceRepository interface {
	SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}
