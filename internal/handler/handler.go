// Package handler provides service handler implementation
package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"code.0x7r.com/scholar7r/harborman/internal/harbor"
	"code.0x7r.com/scholar7r/harborman/internal/notifier"
)

type NotifyHandler struct {
	notifiers []notifier.Notifier
}

func NewNotifyHandler(notifiers []notifier.Notifier) *NotifyHandler {
	return &NotifyHandler{notifiers: notifiers}
}

func (nh *NotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusBadRequest,
		)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(
			"failed to read request body",
			slog.String("error", err.Error()),
		)

		http.Error(
			w,
			"bad request",
			http.StatusBadRequest,
		)

		return
	}

	defer func() { _ = r.Body.Close() }()

	var event harbor.HarborEvent
	if err = json.Unmarshal(body, &event); err != nil {
		slog.Error(
			"failed to unmarshal harbor event",
			slog.String("error", err.Error()),
		)

		http.Error(
			w,
			"bad request",
			http.StatusBadRequest,
		)

		return
	}

	for _, v := range nh.notifiers {
		if err := v.Notify(&event); err != nil {
			slog.Error(
				"failed to send notification",
				slog.String("error", err.Error()),
			)
		}
	}
}
