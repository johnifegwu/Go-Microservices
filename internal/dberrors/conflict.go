package dberrors

type ConflictError struct{}

func (e *ConflictError) Error() string {
	return "attempted to create a record with an existing key"
}
