package lognotify

import "link-society.com/flowg/internal/data/logstorage"

type LogMessage struct {
	Stream   string
	LogKey   string
	LogEntry logstorage.LogEntry
}

type SubscribeMessage struct {
	Stream  string
	SenderC chan<- LogMessage
	DoneC   <-chan struct{}
}
