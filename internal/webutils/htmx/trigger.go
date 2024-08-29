package htmx

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Trigger struct {
	ModalOpenEvent  *ModalOpenEvent  `json:"htmx-custom-modal-open,omitempty"`
	ModalCloseEvent *ModalCloseEvent `json:"htmx-custom-modal-close,omitempty"`
	ToastEvent      *ToastEvent      `json:"htmx-custom-toast,omitempty"`
}

func (t Trigger) Write(ctx context.Context, w http.ResponseWriter) {
	payload, err := json.Marshal(t)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to marshal HX-Trigger payload",
			"channel", "htmx",
			"error", err.Error(),
		)
	} else {
		w.Header().Add("HX-Trigger", string(payload))
	}
}
