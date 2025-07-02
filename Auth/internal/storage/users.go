package storageerrors

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrAlreadyExists    = errors.New("resource already exists")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)
