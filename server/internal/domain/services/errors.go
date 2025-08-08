package services

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrPermissionDenied = errors.New("permission denied")
	ErrConflict         = errors.New("conflict")
)
