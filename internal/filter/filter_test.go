package filter_test

import (
	"reflect"
	"testing"

	"code.0x7r.com/scholar7r/harborman/internal/filter"
	"code.0x7r.com/scholar7r/harborman/internal/harbor"
)

func TestPushEventFilter_Filter(t *testing.T) {
	tests := []struct {
		name string
		opts []filter.PushEventFilterOption
		e    *harbor.Event
		want []harbor.Resource
	}{
		{
			name: "event is nil",
			opts: []filter.PushEventFilterOption{},
			e:    nil,
			want: nil,
		},
		{
			name: "resources is empty",
			opts: []filter.PushEventFilterOption{},
			e: &harbor.Event{
				EventData: harbor.EventData{
					Resources: make([]harbor.Resource, 0),
				},
			},
			want: nil,
		},
		{
			name: "filter no tag",
			opts: []filter.PushEventFilterOption{
				filter.WithFilterNoTag(),
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
			opts: []filter.PushEventFilterOption{
				filter.WithFilterTagEqualDigest(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filter.NewPushEventFilter(tt.opts...)
			got := f.Filter(tt.e)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushEventFilter.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
