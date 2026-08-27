package notifier_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/scholar7r/harborman/v2/internal/notifier"
	"github.com/scholar7r/harborman/v2/pkg/harbor"
)

func TestWebhook_Authorize(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		clientToken string
		want        error
	}{
		{
			name:        "matching token",
			token:       "expected",
			clientToken: "expected",
			want:        nil,
		},
		{
			name:        "mismatched token",
			token:       "expected",
			clientToken: "wrong",
			want:        notifier.ErrTokenMismatch,
		},
		{
			name:        "configured token requires a client token",
			token:       "expected",
			clientToken: "",
			want:        notifier.ErrTokenMismatch,
		},
		{
			name:        "unconfigured token skips verification",
			token:       "",
			clientToken: "anything",
			want:        nil,
		},
		{
			name:        "unconfigured token without client token",
			token:       "",
			clientToken: "",
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := notifier.NewDiscord(
				&http.Client{},
				"https://example.invalid",
				tt.token,
				harbor.NewPushEventFilter(),
			)

			if got := n.Authorize(tt.clientToken); !errors.Is(got, tt.want) {
				t.Errorf("Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhook_Platform(t *testing.T) {
	tests := []struct {
		name string
		n    notifier.Notifier
		want string
	}{
		{
			name: "discord",
			n: notifier.NewDiscord(
				&http.Client{},
				"https://example.invalid",
				"",
				harbor.NewPushEventFilter(),
			),
			want: "discord",
		},
		{
			name: "lark",
			n: notifier.NewLark(
				&http.Client{},
				"https://example.invalid",
				"",
				harbor.NewPushEventFilter(),
			),
			want: "lark",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Platform(); got != tt.want {
				t.Errorf("Platform() = %v, want %v", got, tt.want)
			}
		})
	}
}
