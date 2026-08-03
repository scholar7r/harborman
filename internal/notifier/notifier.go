// Package notifier provides multiple type platform support to send notification
package notifier

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"code.0x7r.com/scholar7r/harborman/internal/cfg"
	"code.0x7r.com/scholar7r/harborman/internal/filter"
	"code.0x7r.com/scholar7r/harborman/internal/harbor"
)

type Notifier interface {
	Notify(context.Context, *harbor.Event) error
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
			// notifiers = append(notifiers, NewDiscord(client, v.URL, pushEventFilter))
		case cfg.NotifierTypeLark:
			notifiers = append(notifiers, NewLark(client, v.URL, pushEventFilter))
		default:
			slog.Warn(
				"unsupported notifier type",
				slog.String("type", string(v.Type)),
			)
		}
	}

	return notifiers
}
