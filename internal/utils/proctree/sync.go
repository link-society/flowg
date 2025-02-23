package proctree

import "sync"

type condValue[T any] struct {
	mu    sync.Mutex
	cond  *sync.Cond
	set   bool
	value T
}

func newCondValue[T any]() *condValue[T] {
	condv := &condValue[T]{}
	condv.cond = sync.NewCond(&condv.mu)
	return condv
}

func (cv *condValue[T]) Wait() T {
	cv.mu.Lock()
	defer cv.mu.Unlock()

	for !cv.set {
		cv.cond.Wait()
	}

	return cv.value
}

func (cv *condValue[T]) Broadcast(value T) {
	cv.set = true
	cv.value = value
	cv.cond.Broadcast()
}
