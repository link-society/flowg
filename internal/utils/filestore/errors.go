package filestore

import "errors"

var (
	ErrStartFailed = errors.New("failed to start filestore")
	ErrStopFailed  = errors.New("failed to stop filestore")
)
