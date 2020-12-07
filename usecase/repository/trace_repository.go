// Package repository contains interfaces to represent functions required to work with storage
package repository

import "github.com/lroman242/redirective/domain"

//go:generate mockgen -package=mocks -destination=mocks/mock_trace_repository.go -source=usecase/repository/trace_repository.go TraceRepository

// TraceRepository interface represent repository to work with trace results.
type TraceRepository interface {
	SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error)
	FindTraceResults(interface{}) (*domain.TraceResults, error)
}
