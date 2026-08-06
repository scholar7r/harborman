package notifier

import (
	"context"
	"net/http"

	"github.com/scholar7r/harborman/internal/filter"
	"github.com/scholar7r/harborman/internal/harbor"
)

type DiscordPayload struct{}

type Discord struct {
	client          *http.Client
	notifyURL       string
	pushEventFilter *filter.PushEventFilter
}

func NewDiscord(client *http.Client, notifyURL string, pushEventFilter *filter.PushEventFilter) Notifier {
	return &Discord{
		client:          client,
		notifyURL:       notifyURL,
		pushEventFilter: pushEventFilter,
	}
}

func (d *Discord) Notify(_ context.Context, _ *harbor.Event) error {
	return nil
}
