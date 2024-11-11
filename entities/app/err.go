package app

import "errors"

var (
	ErrNoRecord             = errors.New("no records found")
	ErrInternalServiceError = errors.New("internal service error")
	ErrDuplicateEntry       = errors.New("duplicate entry")
)
