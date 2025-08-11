package repositories

import "errors"

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrDuplicateISBN    = errors.New("isbn already exists")
	ErrInvalidData      = errors.New("invalid book data")
	ErrInvalidReference = errors.New("invalid reference: the referenced record does not exist")
)
