package dberrors

type ZeroRowsAffectedError struct{}

func (e *ZeroRowsAffectedError) Error() string {
	return "Operation return zero rows affected"
}
