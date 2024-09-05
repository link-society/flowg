package webutils

import (
	"context"
)

type notifySystem struct {
	notifications []string
}

type notifyContextKey string

const CONTEXT_NOTIFY_SYSTEM notifyContextKey = "notify_system"

func WithNotificationSystem(ctx context.Context) context.Context {
	return context.WithValue(ctx, CONTEXT_NOTIFY_SYSTEM, &notifySystem{
		notifications: make([]string, 0),
	})
}

func NotifyInfo(ctx context.Context, message string) {
	notify(ctx, "&#9989; "+message)
}

func NotifyWarning(ctx context.Context, message string) {
	notify(ctx, "&#9888;&#65039; "+message)
}

func NotifyError(ctx context.Context, message string) {
	notify(ctx, "&#10060; "+message)
}

func notify(ctx context.Context, message string) {
	system := ctx.Value(CONTEXT_NOTIFY_SYSTEM).(*notifySystem)
	system.notifications = append(system.notifications, message)
}

func Notifications(ctx context.Context) []string {
	system := ctx.Value(CONTEXT_NOTIFY_SYSTEM)
	if system == nil {
		return nil
	}

	return system.(*notifySystem).notifications
}
