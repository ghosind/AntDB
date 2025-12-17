package core

import "errors"

var (
	ErrWrongType  = errors.New("wrong type")
	ErrNoSuchKey  = errors.New("no such key")
	ErrNotInteger = errors.New("value is not an integer or out of range")
)
