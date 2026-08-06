package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/scholar7r/harborman/internal/filter"
	"github.com/scholar7r/harborman/internal/harbor"
)

type LarkPayload struct {
	MsgType string   `json:"msg_type"` // interactive
	Card    LarkCard `json:"card"`
}

type LarkCard struct {
	Schema string `json:"schema"` // "2.0"
	Header any    `json:"header"`
	Body   any    `json:"body"`
}

type LarkCardHeader struct {
	Template string    `json:"template"`
	Title    LarkField `json:"title"`
}

type LarkCardBody struct {
	Elements []LarkField `json:"elements"`
}

type LarkField struct {
	Tag      string      `json:"tag"`
	Content  string      `json:"content,omitempty"`
	Elements []LarkField `json:"elements,omitempty"`
}

type Lark struct {
	client          *http.Client
	notifyURL       string
	pushEventFilter *filter.PushEventFilter
}

func NewLark(client *http.Client, notifyURL string, pushEventFilter *filter.PushEventFilter) Notifier {
	return &Lark{
		client:          client,
		notifyURL:       notifyURL,
		pushEventFilter: pushEventFilter,
	}
}

func (l *Lark) Notify(ctx context.Context, e *harbor.Event) error {
	var payload *LarkPayload
	switch e.Type {
	case harbor.EventPushArtifact:
		e.EventData.Resources = l.pushEventFilter.Filter(e)

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

		payload = &LarkPayload{
			MsgType: "interactive",
			Card: LarkCard{
				Schema: "2.0",
				Header: LarkCardHeader{
					Template: "blue",
					Title: LarkField{
						Tag:     "plain_text",
						Content: "Harbor Push Notification",
					},
				},
				Body: LarkCardBody{
					Elements: []LarkField{
						{
							Tag:     "markdown",
							Content: fmt.Sprintf("Repository: %s", e.EventData.Repository.Name),
						},
						{
							Tag:     "markdown",
							Content: fmt.Sprintf("Tags: %s", strings.Join(tags, ", ")),
						},
					},
				},
			},
		}

	default:
		return fmt.Errorf("unsupported event type: %s", e.Type)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return l.sendNotification(ctx, data)
}

func (l *Lark) sendNotification(ctx context.Context, b []byte) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		l.notifyURL,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read lark response: %w", err)
	}

	slog.DebugContext(
		ctx, "lark response",
		slog.Int("status_code", resp.StatusCode),
		slog.String("body", string(respBody)),
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lark returned abnormal status code %d", resp.StatusCode)
	}

	return nil
}
