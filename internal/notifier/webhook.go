package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/scholar7r/harborman/internal/filter"
	"github.com/scholar7r/harborman/internal/harbor"
)

type payloadBuilder interface {
	Platform() string
	BuildPushArtifact(e *harbor.Event, tags []string) any
}

type webhook struct {
	client          *http.Client
	notifyURL       string
	pushEventFilter *filter.PushEventFilter
	builder         payloadBuilder
}

func newWebhook(
	client *http.Client,
	notifyURL string,
	pushEventFilter *filter.PushEventFilter,
	builder payloadBuilder,
) Notifier {
	return &webhook{
		client:          client,
		notifyURL:       notifyURL,
		pushEventFilter: pushEventFilter,
		builder:         builder,
	}
}

func (w *webhook) Notify(ctx context.Context, e *harbor.Event) error {
	var payload any

	switch e.Type {
	case harbor.EventPushArtifact:
		e.EventData.Resources = w.pushEventFilter.Filter(e)

		tags := make([]string, 0, len(e.EventData.Resources))
		for _, v := range e.EventData.Resources {
			tags = append(tags, v.Tag)
		}

		if len(tags) == 0 {
			slog.WarnContext(
				ctx,
				"no valid resources to notify after filtering",
				slog.String("repository", e.EventData.Repository.Name),
			)
		}

		payload = w.builder.BuildPushArtifact(e, tags)

	default:
		return fmt.Errorf("unsupported event type: %s", e.Type)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return w.send(ctx, data)
}

func (w *webhook) send(ctx context.Context, b []byte) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		w.notifyURL,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	slog.DebugContext(
		ctx, "notifier response",
		slog.String("platform", w.builder.Platform()),
		slog.Int("status_code", resp.StatusCode),
		slog.String("body", string(respBody)),
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"%s returned abnormal status code %d",
			w.builder.Platform(), resp.StatusCode,
		)
	}

	return nil
}
