// Package handler provides service handler implementation
package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/scholar7r/harborman/internal/cfg"
	"github.com/scholar7r/harborman/internal/harbor"
	"github.com/scholar7r/harborman/internal/notifier"
)

type NotifyHandler struct {
	c         *cfg.Cfg
	notifiers []notifier.Notifier
}

func NewNotifyHandler(c *cfg.Cfg, notifiers []notifier.Notifier) *NotifyHandler {
	return &NotifyHandler{
		c:         c,
		notifiers: notifiers,
	}
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

	wg := &sync.WaitGroup{}
	clientToken := r.Header.Get("Authorization")

	for i, v := range nh.notifiers {
		wg.Add(1)
		go func(idx int, n notifier.Notifier) {
			defer wg.Done()

			if clientToken != "" {
				localToken := nh.c.Notifiers[idx].Authorization
				if localToken == "" {
					localToken = nh.c.Authorization
				}

				if localToken != "" && localToken != clientToken {
					slog.ErrorContext(
						r.Context(),
						"incorrect authorization token",
						slog.Int("notifier_index", idx),
					)
					return
				}
			}

			if ce := n.Notify(r.Context(), &event); ce != nil {
				slog.ErrorContext(
					r.Context(),
					"failed to send notification",
					slog.String("error", ce.Error()),
				)
			}
		}(i, v)
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)
}
