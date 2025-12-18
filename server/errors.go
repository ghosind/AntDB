package server

import "errors"

var (
	ErrSyntax          = errors.New("syntax error")
	ErrInvalidDBIndex  = errors.New("value is not an integer or out of range")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNotPermitted    = errors.New("operation not permitted")
)

func newUnknownCommandError(cmd string) error {
	return errors.New("unknown command '" + cmd + "'")
}

func newWrongArityError(cmd string) error {
	return errors.New("wrong number of arguments for '" + cmd + "' command")
}
