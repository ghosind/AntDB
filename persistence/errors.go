package persistence

import "errors"

var (
	ErrInvalidRDBFormat      = errors.New("invalid RDB format")
	ErrUnsupportedObjectType = errors.New("unsupported object type")
)
