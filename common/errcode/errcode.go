package errcode

import "errors"

var (
	ErrNoItems = errors.New("items cannot be empty")
)
