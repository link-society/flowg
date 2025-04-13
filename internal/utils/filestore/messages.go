package filestore

import (
	"io/fs"
	"strings"
)

type message interface {
	Handle(*procHandler)
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

var _ message = (*listItems)(nil)
var _ message = (*readItem)(nil)
var _ message = (*writeItem)(nil)
var _ message = (*deleteItem)(nil)
var _ message = (*statItem)(nil)

func (msg *listItems) Handle(handler *procHandler) {
	keys := []string{}

	handler.cache.Range(func(key, _ interface{}) bool {
		filename := key.(string)

		if handler.extension == "" || strings.HasSuffix(filename, handler.extension) {
			suffixLen := len(handler.extension)
			itemName := filename[:len(filename)-suffixLen]
			keys = append(keys, itemName)
		}

		return true
	})

	msg.replyTo.SendOk(keys)
}

func (msg *readItem) Handle(handler *procHandler) {
	filename := msg.key + handler.extension

	value, ok := handler.cache.Load(filename)
	if !ok {
		handler.mu.Lock()
		content, err := handler.backend.ReadFile(filename)
		handler.mu.Unlock()

		if err != nil {
			msg.replyTo.SendErr(err)
		} else {
			handler.cache.Store(filename, content)
			msg.replyTo.SendOk(content)
		}
	} else {
		msg.replyTo.SendOk(value.([]byte))
	}
}

func (msg *writeItem) Handle(handler *procHandler) {
	filename := msg.key + handler.extension

	handler.mu.Lock()
	err := handler.backend.WriteFile(filename, msg.content)
	handler.mu.Unlock()

	if err != nil {
		msg.replyTo.SendErr(err)
	} else {
		handler.cache.Store(filename, msg.content)
		msg.replyTo.SendOk(true)
	}
}

func (msg *deleteItem) Handle(handler *procHandler) {
	filename := msg.key + handler.extension

	handler.mu.Lock()
	err := handler.backend.DeleteFile(filename)
	handler.mu.Unlock()

	if err != nil {
		msg.replyTo.SendErr(err)
	} else {
		handler.cache.Delete(filename)
		msg.replyTo.SendOk(true)
	}
}

func (msg *statItem) Handle(handler *procHandler) {
	filename := msg.key + handler.extension

	handler.mu.Lock()
	info, err := handler.backend.StatFile(filename)
	handler.mu.Unlock()

	if err != nil {
		msg.replyTo.SendErr(err)
	} else {
		msg.replyTo.SendOk(statItemResponse{info: info})
	}
}
