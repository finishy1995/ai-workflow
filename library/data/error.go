package data

import "errors"

var (
	ErrMustInit             = errors.New("must init first")
	ErrParamInvalid         = errors.New("config must be struct pointer")
	ErrCannotCreateInstance = errors.New("cannot create object instance")
)
