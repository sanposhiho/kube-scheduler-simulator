package errors

import "errors"

var (
	ErrTooManyRunningScenario = errors.New("too many running scenario")
	ErrNoRunningScenario      = errors.New("no running scenario")
)
