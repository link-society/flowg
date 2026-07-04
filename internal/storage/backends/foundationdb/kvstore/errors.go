package kvstore

import "errors"

var (
	ErrStartFailed = errors.New("failed to start foundationdb kvstore")
	ErrStopFailed  = errors.New("failed to stop foundationdb kvstore")
)
