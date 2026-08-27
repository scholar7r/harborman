package harbor_test

import (
	"reflect"
	"testing"

	"github.com/scholar7r/harborman/pkg/harbor"
)

func TestPushEventFilter_Filter(t *testing.T) {
	tests := []struct {
		name string
		opts []harbor.PushEventFilterOption
		e    *harbor.Event
		want []harbor.Resource
	}{
		{
			name: "event is nil",
			opts: []harbor.PushEventFilterOption{},
			e:    nil,
			want: nil,
		},
		{
			name: "resources is empty",
			opts: []harbor.PushEventFilterOption{},
			e: &harbor.Event{
				EventData: harbor.EventData{
					Resources: make([]harbor.Resource, 0),
				},
			},
			want: nil,
		},
		{
			name: "filter no tag",
			opts: []harbor.PushEventFilterOption{
				harbor.WithFilterNoTag(),
			},
			e: &harbor.Event{
				EventData: harbor.EventData{
					Resources: []harbor.Resource{
						{Tag: "", Digest: "sha256:f1f72d58-224d-4263-8d49-6bcb517c27c9"}, // filtered
						{Tag: "v1.0.0", Digest: "sha256:ef3891a0-8e24-45ae-a163-0504b6bfea9c"},
					},
				},
			},
			want: []harbor.Resource{
				{Tag: "v1.0.0", Digest: "sha256:ef3891a0-8e24-45ae-a163-0504b6bfea9c"},
			},
		},
		{
			name: "filter tag equal digest",
			opts: []harbor.PushEventFilterOption{
				harbor.WithFilterTagEqualDigest(),
			},
			e: &harbor.Event{
				EventData: harbor.EventData{
					Resources: []harbor.Resource{
						{
							Tag:    "sha256:f1f72d58-224d-4263-8d49-6bcb517c27c9",
							Digest: "sha256:f1f72d58-224d-4263-8d49-6bcb517c27c9",
						}, // filtered
						{Tag: "v1.0.0", Digest: "sha256:ef3891a0-8e24-45ae-a163-0504b6bfea9c"},
					},
				},
			},
			want: []harbor.Resource{
				{Tag: "v1.0.0", Digest: "sha256:ef3891a0-8e24-45ae-a163-0504b6bfea9c"},
			},
		},
		{
			name: "filter no tag and tag equal digest",
			opts: []harbor.PushEventFilterOption{
				harbor.WithFilterNoTag(),
				harbor.WithFilterTagEqualDigest(),
			},
			e: &harbor.Event{
				EventData: harbor.EventData{
					Resources: []harbor.Resource{
						{
							Tag:    "sha256:30dda09d-54e3-4682-8cee-28ea282bdcae",
							Digest: "sha256:30dda09d-54e3-4682-8cee-28ea282bdcae",
						}, // filtered
						{Tag: "", Digest: "sha256:fe62911f-bba0-43ed-ac47-826372ea18c4"}, // filtered
						{Tag: "v1.0.0", Digest: "sha256:65879c47-d731-4d84-99b1-e1b635e2f78c"},
					},
				},
			},
			want: []harbor.Resource{
				{Tag: "v1.0.0", Digest: "sha256:65879c47-d731-4d84-99b1-e1b635e2f78c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := harbor.NewPushEventFilter(tt.opts...)
			got := f.Filter(tt.e)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushEventFilter.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
