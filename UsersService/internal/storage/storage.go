package storageerror

import "errors"

var (
	ErrNotFound      = errors.New("resourse not found")
	ErrAlreadyExists = errors.New("resourse already exists")
)
