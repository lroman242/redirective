package storage

//InvalidIDTypeError describe error that occurs when wrong ID type received
type InvalidIDTypeError struct {
}

//Error function return error message
func (e *InvalidIDTypeError) Error() string {
	return "invalid id type received"
}
