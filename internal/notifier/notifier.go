// Package notifier provides multiple type platform support to send notification
package notifier

import (
	"code.0x7r.com/scholar7r/harborman/internal/cfg"
	"code.0x7r.com/scholar7r/harborman/internal/harbor"
)

type Notifier interface {
	Type() string
	Notify(*harbor.HarborEvent) error
}

func FromCfg(nc []cfg.NotifierCfg) []Notifier {
	var notifiers []Notifier
	for _, v := range nc {
		switch v.Type {
		case cfg.NotifierTypeDiscord:
			notifiers = append(notifiers, NewDiscordNotifier(v.URL))
		case cfg.NotifierTypeLark:
			notifiers = append(notifiers, NewLarkNotifier(v.URL))
		}
	}

	return notifiers
}
