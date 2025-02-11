package filestore

import (
	"io/fs"
	"strings"
)

type message interface {
	Handle(*workerRunning)
}

type replyTo[T any] struct {
	okC  chan T
	errC chan error
}

func (r replyTo[T]) SendOk(value T) {
	r.okC <- value
	close(r.okC)
	close(r.errC)
}

func (r replyTo[T]) SendErr(err error) {
	r.errC <- err
	close(r.okC)
	close(r.errC)
}

func (r replyTo[T]) Receive() (T, error) {
	select {
	case value := <-r.okC:
		return value, nil
	case err := <-r.errC:
		var zero T
		return zero, err
	}
}

type listItems struct {
	replyTo replyTo[[]string]
}

type readItem struct {
	replyTo replyTo[[]byte]
	key     string
}

type writeItem struct {
	replyTo replyTo[bool]
	key     string
	content []byte
}

type deleteItem struct {
	replyTo replyTo[bool]
	key     string
}

type statItemResponse struct {
	info fs.FileInfo
}

type statItem struct {
	replyTo replyTo[statItemResponse]
	key     string
}

func (msg *listItems) Handle(workerState *workerRunning) {
	keys := []string{}

	workerState.cache.Range(func(key, _ interface{}) bool {
		filename := key.(string)

		if workerState.extension == "" || strings.HasSuffix(filename, workerState.extension) {
			suffixLen := len(workerState.extension)
			itemName := filename[:len(filename)-suffixLen]
			keys = append(keys, itemName)
		}

		return true
	})

	msg.replyTo.SendOk(keys)
}

func (msg *readItem) Handle(workerState *workerRunning) {
	filename := msg.key + workerState.extension

	value, ok := workerState.cache.Load(filename)
	if !ok {
		workerState.mu.Lock()
		content, err := workerState.backend.ReadFile(filename)
		workerState.mu.Unlock()

		if err != nil {
			msg.replyTo.SendErr(err)
		} else {
			workerState.cache.Store(filename, content)
			msg.replyTo.SendOk(content)
		}
	} else {
		msg.replyTo.SendOk(value.([]byte))
	}
}

func (msg *writeItem) Handle(workerState *workerRunning) {
	filename := msg.key + workerState.extension

	workerState.mu.Lock()
	err := workerState.backend.WriteFile(filename, msg.content)
	workerState.mu.Unlock()

	if err != nil {
		msg.replyTo.SendErr(err)
	} else {
		workerState.cache.Store(filename, msg.content)
		msg.replyTo.SendOk(true)
	}
}

func (msg *deleteItem) Handle(workerState *workerRunning) {
	filename := msg.key + workerState.extension

	workerState.mu.Lock()
	err := workerState.backend.DeleteFile(filename)
	workerState.mu.Unlock()

	if err != nil {
		msg.replyTo.SendErr(err)
	} else {
		workerState.cache.Delete(filename)
		msg.replyTo.SendOk(true)
	}
}

func (msg *statItem) Handle(workerState *workerRunning) {
	filename := msg.key + workerState.extension

	workerState.mu.Lock()
	info, err := workerState.backend.StatFile(filename)
	workerState.mu.Unlock()

	if err != nil {
		msg.replyTo.SendErr(err)
	} else {
		msg.replyTo.SendOk(statItemResponse{info: info})
	}
}
