package notifier_test

import (
	"reflect"
	"testing"

	"github.com/scholar7r/harborman/v2/internal/notifier"
	"github.com/scholar7r/harborman/v2/pkg/harbor"
)

func TestDiscord_BuildPushArtifact(t *testing.T) {
	tests := []struct {
		name string
		e    *harbor.Event
		tags []string
		want any
	}{
		{
			name: "single tag event",
			e: &harbor.Event{
				Type:     harbor.EventPushArtifact,
				OccurAt:  1754438400,
				Operator: "scholar7r",
				EventData: harbor.EventData{
					Repository: harbor.Repository{
						Name:      "harborman",
						Namespace: "scholar7r",
						FullName:  "scholar7r/harborman",
					},
				},
			},
			tags: []string{"1.0.0"},
			want: &notifier.DiscordPayload{
				Embeds: []notifier.DiscordEmbed{
					{
						Title: "Harbor Push Notification",
						Color: 0x57F287,
						Fields: []notifier.EmbedField{
							{Name: "Project", Value: "scholar7r", Inline: true},
							{Name: "Repository", Value: "scholar7r/harborman", Inline: true},
							{Name: "Pushed by", Value: "scholar7r", Inline: true},
							{Name: "Tags", Value: "1.0.0", Inline: true},
						},
						Timestamp: "2025-08-06T00:00:00Z",
						Footer:    &notifier.EmbedFooter{Text: "Harbor Registry"},
					},
				},
			},
		},
		{
			name: "multiple tags event",
			e: &harbor.Event{
				Type:     harbor.EventPushArtifact,
				OccurAt:  1754438400,
				Operator: "scholar7r",
				EventData: harbor.EventData{
					Repository: harbor.Repository{
						Name:      "harborman",
						Namespace: "scholar7r",
						FullName:  "scholar7r/harborman",
					},
				},
			},
			tags: []string{"latest", "1.0.0", "1.0"},
			want: &notifier.DiscordPayload{
				Embeds: []notifier.DiscordEmbed{
					{
						Title: "Harbor Push Notification",
						Color: 0x57F287,
						Fields: []notifier.EmbedField{
							{Name: "Project", Value: "scholar7r", Inline: true},
							{Name: "Repository", Value: "scholar7r/harborman", Inline: true},
							{Name: "Pushed by", Value: "scholar7r", Inline: true},
							{Name: "Tags", Value: "latest, 1.0.0, 1.0", Inline: true},
						},
						Timestamp: "2025-08-06T00:00:00Z",
						Footer:    &notifier.EmbedFooter{Text: "Harbor Registry"},
					},
				},
			},
		},
		{
			name: "no tag",
			e: &harbor.Event{
				Type:     harbor.EventPushArtifact,
				OccurAt:  1754438400,
				Operator: "scholar7r",
				EventData: harbor.EventData{
					Repository: harbor.Repository{
						Name:      "harborman",
						Namespace: "scholar7r",
						FullName:  "scholar7r/harborman",
					},
				},
			},
			tags: []string{},
			want: &notifier.DiscordPayload{
				Embeds: []notifier.DiscordEmbed{
					{
						Title: "Harbor Push Notification",
						Color: 0x57F287,
						Fields: []notifier.EmbedField{
							{Name: "Project", Value: "scholar7r", Inline: true},
							{Name: "Repository", Value: "scholar7r/harborman", Inline: true},
							{Name: "Pushed by", Value: "scholar7r", Inline: true},
							{Name: "Tags", Value: "", Inline: true},
						},
						Timestamp: "2025-08-06T00:00:00Z",
						Footer:    &notifier.EmbedFooter{Text: "Harbor Registry"},
					},
				},
			},
		},
		{
			name: "nil tag",
			e: &harbor.Event{
				Type:     harbor.EventPushArtifact,
				OccurAt:  1754438400,
				Operator: "scholar7r",
				EventData: harbor.EventData{
					Repository: harbor.Repository{
						Name:      "harborman",
						Namespace: "scholar7r",
						FullName:  "scholar7r/harborman",
					},
				},
			},
			tags: nil,
			want: &notifier.DiscordPayload{
				Embeds: []notifier.DiscordEmbed{
					{
						Title: "Harbor Push Notification",
						Color: 0x57F287,
						Fields: []notifier.EmbedField{
							{Name: "Project", Value: "scholar7r", Inline: true},
							{Name: "Repository", Value: "scholar7r/harborman", Inline: true},
							{Name: "Pushed by", Value: "scholar7r", Inline: true},
							{Name: "Tags", Value: "", Inline: true},
						},
						Timestamp: "2025-08-06T00:00:00Z",
						Footer:    &notifier.EmbedFooter{Text: "Harbor Registry"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d notifier.Discord
			got := d.BuildPushArtifact(tt.e, tt.tags)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildPushArtifact() = %v, want %v", got, tt.want)
			}
		})
	}
}
