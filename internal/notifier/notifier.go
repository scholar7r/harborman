// Package notifier provides multiple type platform support to send notification
package notifier

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/scholar7r/harborman/internal/cfg"
	"github.com/scholar7r/harborman/internal/filter"
	"github.com/scholar7r/harborman/internal/harbor"
)

var ErrTokenMismatch = errors.New("authorization token mismatch")

type Notifier interface {
	Platform() string
	Authorize(clientToken string) error
	Notify(ctx context.Context, e *harbor.Event) error
}

func FromCfg(client *http.Client, nc []cfg.NotifierCfg) []Notifier {
	var notifiers []Notifier

	for i, v := range nc {
		if v.Type == "" || v.URL == "" {
			slog.Warn(
				"skip notifier with incomplete fields",
				slog.Int("index", i),
			)

			continue
		}

		_, err := url.Parse(v.URL)
		if err != nil {
			slog.Error(
				"skip notifier with invalid URL",
				slog.Int("index", i),
				slog.String("url", v.URL),
			)

			continue
		}

		pushEventFilter := filter.NewPushEventFilter(
			filter.WithFilterNoTag(),
			filter.WithFilterTagEqualDigest(),
		)

		switch v.Type {
		case cfg.NotifierTypeDiscord:
			notifiers = append(
				notifiers,
				NewDiscord(client, v.URL, v.Authorization, pushEventFilter),
			)
		case cfg.NotifierTypeLark:
			notifiers = append(
				notifiers,
				NewLark(client, v.URL, v.Authorization, pushEventFilter),
			)
		default:
			slog.Warn(
				"unsupported notifier type",
				slog.String("type", string(v.Type)),
			)
		}
	}

	return notifiers
}
