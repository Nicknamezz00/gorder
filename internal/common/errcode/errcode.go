package errcode

import "errors"

var (
	ErrNoItems       = errors.New("items cannot be empty")
	ErrOrderNotFound = errors.New("order not found")
)
