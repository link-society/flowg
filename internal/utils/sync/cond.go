package sync

import (
	gosync "sync"
)

type CondValue[T any] struct {
	mu    gosync.Mutex
	cond  *gosync.Cond
	set   bool
	value T
}

func NewCondValue[T any]() *CondValue[T] {
	condv := &CondValue[T]{}
	condv.cond = gosync.NewCond(&condv.mu)
	return condv
}

func (cv *CondValue[T]) Wait() T {
	cv.mu.Lock()
	defer cv.mu.Unlock()

	for !cv.set {
		cv.cond.Wait()
	}

	return cv.value
}

func (cv *CondValue[T]) Broadcast(value T) {
	cv.set = true
	cv.value = value
	cv.cond.Broadcast()
}
