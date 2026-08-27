package notifier

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/scholar7r/harborman/internal/cfg"
	"github.com/scholar7r/harborman/pkg/harbor"
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

type Lark struct{}

func NewLark(
	client *http.Client,
	notifyURL string,
	token string,
	pushEventFilter *harbor.PushEventFilter,
) Notifier {
	return newWebhook(client, notifyURL, token, pushEventFilter, &Lark{})
}

func (l *Lark) Platform() string {
	return string(cfg.NotifierTypeLark)
}

func (l *Lark) BuildPushArtifact(e *harbor.Event, tags []string) any {
	return &LarkPayload{
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
}
