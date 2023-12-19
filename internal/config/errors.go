package config

import "errors"

var (
	ErrRequiredVariable = errors.New("missing required variable")
	ErrEmptyVariable    = errors.New("empty variable")
	ErrUnsupportedType  = errors.New("unsupported type")
	ErrUnknownParam     = errors.New("unknown param")
)
