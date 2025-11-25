package database

import "errors"

var (
	ErrWrongType = errors.New("wrong type")
	ErrNoSuchKey = errors.New("no such key")
)
