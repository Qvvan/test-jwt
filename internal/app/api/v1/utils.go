package v1

type PublicError struct {
	err    error
	status int
}

func NewPublicErr(err error, status int) *PublicError {
	return &PublicError{err: err, status: status}
}

func (e *PublicError) Error() string {
	return e.err.Error()
}
