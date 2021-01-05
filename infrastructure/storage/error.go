package storage

import "fmt"

// InvalidIDTypeError describe error that occurs when wrong ID type received.
type InvalidIDTypeError struct {
}

// Error function return error message.
func (e *InvalidIDTypeError) Error() string {
	return "invalid id type received"
}

// CannotDecodeRecordError describe error that occurs when response from db is not decoded into provided struct.
type CannotDecodeRecordError struct {
	err error
}

// Error function return error message.
func (e *CannotDecodeRecordError) Error() string {
	return fmt.Sprintf("mongodb reulsts decode failed. error: %s", e.err)
}
