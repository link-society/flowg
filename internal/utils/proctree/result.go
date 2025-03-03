package proctree

import (
	"context"
	"errors"

	"github.com/vladopajic/go-actor/actor"
)

type ProcessResult interface {
	Error() error
	Done() bool
}

func Continue() ProcessResult {
	return &continueR{}
}

func Terminate(err error) ProcessResult {
	switch {
	case errors.Is(err, context.Canceled):
		return &terminateR{err: nil}

	case errors.Is(err, actor.ErrStopped):
		return &terminateR{err: nil}

	default:
		return &terminateR{err: err}
	}
}

type continueR struct{}

func (r *continueR) Error() error {
	return nil
}

func (r *continueR) Done() bool {
	return false
}

type terminateR struct {
	err error
}

func (r *terminateR) Error() error {
	return r.err
}

func (r *terminateR) Done() bool {
	return true
}
