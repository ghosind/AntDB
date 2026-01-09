package client

import "errors"

var (
	ErrEmptyCommand      = errors.New("empty command")
	ErrInvalidArray      = errors.New("invalid array format")
	ErrInvalidBulkHeader = errors.New("invalid bulk header")
)
