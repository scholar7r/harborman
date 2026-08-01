// Package notifier provides multiple type platform support to send notification
package notifier

import (
	"log/slog"

	"code.0x7r.com/scholar7r/harborman/internal/cfg"
	"code.0x7r.com/scholar7r/harborman/internal/harbor"
)

type Notifier interface {
	Notify(*harbor.Event) error
}

func FromCfg(nc []cfg.NotifierCfg) []Notifier {
	var notifiers []Notifier
	for i, v := range nc {
		if v.Type == "" || v.URL == "" {
			slog.Warn(
				"skip notifier with incomplete fields",
				slog.Int("index", i),
			)

			continue
		}

		switch v.Type {
		case cfg.NotifierTypeDiscord:
			// notifiers = append(notifiers, NewDiscord(v.URL))
		case cfg.NotifierTypeLark:
			// notifiers = append(notifiers, NewLark(v.URL))
		default:
			slog.Warn(
				"unsupported notifier type",
				slog.String("type", string(v.Type)),
			)
		}
	}

	return notifiers
}
