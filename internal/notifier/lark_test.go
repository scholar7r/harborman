package notifier_test

import (
	"reflect"
	"testing"

	"github.com/scholar7r/harborman/internal/notifier"
	"github.com/scholar7r/harborman/pkg/harbor"
)

func TestLark_BuildPushArtifact(t *testing.T) {
	tests := []struct {
		name string
		e    *harbor.Event
		tags []string
		want any
	}{
		{
			name: "single tag event",
			e: &harbor.Event{
				Type: harbor.EventPushArtifact,
				EventData: harbor.EventData{
					Repository: harbor.Repository{Name: "scholar7r/harborman"},
				},
			},
			tags: []string{"1.0.0"},
			want: &notifier.LarkPayload{
				MsgType: "interactive",
				Card: notifier.LarkCard{
					Schema: "2.0",
					Header: notifier.LarkCardHeader{
						Template: "blue",
						Title: notifier.LarkField{
							Tag:     "plain_text",
							Content: "Harbor Push Notification",
						},
					},
					Body: notifier.LarkCardBody{
						Elements: []notifier.LarkField{
							{Tag: "markdown", Content: "Repository: scholar7r/harborman"},
							{Tag: "markdown", Content: "Tags: 1.0.0"},
						},
					},
				},
			},
		},
		{
			name: "multiple tags event",
			e: &harbor.Event{
				Type: harbor.EventPushArtifact,
				EventData: harbor.EventData{
					Repository: harbor.Repository{Name: "scholar7r/harborman"},
				},
			},
			tags: []string{"latest", "1.0.0", "1.0"},
			want: &notifier.LarkPayload{
				MsgType: "interactive",
				Card: notifier.LarkCard{
					Schema: "2.0",
					Header: notifier.LarkCardHeader{
						Template: "blue",
						Title: notifier.LarkField{
							Tag:     "plain_text",
							Content: "Harbor Push Notification",
						},
					},
					Body: notifier.LarkCardBody{
						Elements: []notifier.LarkField{
							{Tag: "markdown", Content: "Repository: scholar7r/harborman"},
							{Tag: "markdown", Content: "Tags: latest, 1.0.0, 1.0"},
						},
					},
				},
			},
		},
		{
			name: "no tag",
			e: &harbor.Event{
				Type: harbor.EventPushArtifact,
				EventData: harbor.EventData{
					Repository: harbor.Repository{Name: "scholar7r/harborman"},
				},
			},
			tags: []string{},
			want: &notifier.LarkPayload{
				MsgType: "interactive",
				Card: notifier.LarkCard{
					Schema: "2.0",
					Header: notifier.LarkCardHeader{
						Template: "blue",
						Title: notifier.LarkField{
							Tag:     "plain_text",
							Content: "Harbor Push Notification",
						},
					},
					Body: notifier.LarkCardBody{
						Elements: []notifier.LarkField{
							{Tag: "markdown", Content: "Repository: scholar7r/harborman"},
							{Tag: "markdown", Content: "Tags: "},
						},
					},
				},
			},
		},
		{
			name: "nil tag",
			e: &harbor.Event{
				Type: harbor.EventPushArtifact,
				EventData: harbor.EventData{
					Repository: harbor.Repository{Name: "scholar7r/harborman"},
				},
			},
			tags: nil,
			want: &notifier.LarkPayload{
				MsgType: "interactive",
				Card: notifier.LarkCard{
					Schema: "2.0",
					Header: notifier.LarkCardHeader{
						Template: "blue",
						Title: notifier.LarkField{
							Tag:     "plain_text",
							Content: "Harbor Push Notification",
						},
					},
					Body: notifier.LarkCardBody{
						Elements: []notifier.LarkField{
							{Tag: "markdown", Content: "Repository: scholar7r/harborman"},
							{Tag: "markdown", Content: "Tags: "},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l notifier.Lark
			got := l.BuildPushArtifact(tt.e, tt.tags)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildPushArtifact() = %v, want %v", got, tt.want)
			}
		})
	}
}
