package kvstore

import "errors"

var (
	ErrStartFailed = errors.New("failed to start kvstore")
	ErrStopFailed  = errors.New("failed to stop kvstore")
)
