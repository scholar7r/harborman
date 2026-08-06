package notifier

import (
	"net/http"
	"strings"
	"time"

	"github.com/scholar7r/harborman/internal/cfg"
	"github.com/scholar7r/harborman/internal/filter"
	"github.com/scholar7r/harborman/internal/harbor"
)

const discordColorGreen = 0x57F287

type DiscordPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Timestamp   string       `json:"timestamp,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type Discord struct{}

func NewDiscord(client *http.Client, notifyURL string, pushEventFilter *filter.PushEventFilter) Notifier {
	return newWebhook(client, notifyURL, pushEventFilter, &Discord{})
}

func (d *Discord) Platform() string {
	return string(cfg.NotifierTypeDiscord)
}

func (d *Discord) BuildPushArtifact(e *harbor.Event, tags []string) any {
	return &DiscordPayload{
		Embeds: []DiscordEmbed{
			{
				Title: "Harbor Push Notification",
				Color: discordColorGreen,
				Fields: []EmbedField{
					{
						Name:   "Project",
						Value:  e.EventData.Repository.Namespace,
						Inline: true,
					},
					{
						Name:   "Repository",
						Value:  e.EventData.Repository.FullName,
						Inline: true,
					},
					{
						Name:   "Pushed by",
						Value:  e.Operator,
						Inline: true,
					},
					{
						Name:   "Tags",
						Value:  strings.Join(tags, ", "),
						Inline: true,
					},
				},
				Timestamp: time.Unix(e.OccurAt, 0).UTC().Format(time.RFC3339),
				Footer:    &EmbedFooter{Text: "Harbor Registry"},
			},
		},
	}
}
