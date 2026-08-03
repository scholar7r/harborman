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
	slog.InfoContext(
		r.Context(),
		"income request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("content_type", r.Header.Get("Content-Type")),
		slog.Int64("content_length", r.ContentLength),
	)

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(
			r.Context(),
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

	var event harbor.Event
	if err = json.Unmarshal(body, &event); err != nil {
		slog.ErrorContext(
			r.Context(),
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
		if err = v.Notify(r.Context(), &event); err != nil {
			slog.ErrorContext(
				r.Context(),
				"failed to send notification",
				slog.String("error", err.Error()),
			)
		}
	}
}
