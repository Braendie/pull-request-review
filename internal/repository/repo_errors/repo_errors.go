package repoerrors

import "errors"

var (
	ErrAlreadyExistsCode = "unique_violation"
	ErrAlreadyExists     = errors.New("Already exists")
	ErrNotFound = errors.New("Not Found")
)
