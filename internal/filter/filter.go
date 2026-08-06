// Package filter provides the ability to filter harbor events
package filter

import "github.com/scholar7r/harborman/internal/harbor"

type (
	PushEventFilterOption func(*PushEventFilter)
	PushEventFilter       struct {
		filterNoTag          bool
		filterTagEqualDigest bool
	}
)

func NewPushEventFilter(opts ...PushEventFilterOption) *PushEventFilter {
	pushEventFilter := &PushEventFilter{}

	for _, opt := range opts {
		opt(pushEventFilter)
	}

	return pushEventFilter
}

func (f *PushEventFilter) Filter(e *harbor.Event) []harbor.Resource {
	if e == nil || len(e.EventData.Resources) == 0 {
		return nil
	}

	var filtered []harbor.Resource

	for _, v := range e.EventData.Resources {
		if f.filterNoTag && isNoTag(v) {
			continue
		}

		if f.filterTagEqualDigest && isTagEqualDigest(v) {
			continue
		}

		filtered = append(filtered, v)
	}

	return filtered
}

func WithFilterNoTag() PushEventFilterOption {
	return func(f *PushEventFilter) {
		f.filterNoTag = true
	}
}

func WithFilterTagEqualDigest() PushEventFilterOption {
	return func(f *PushEventFilter) {
		f.filterTagEqualDigest = true
	}
}

func isNoTag(r harbor.Resource) bool {
	return r.Tag == ""
}

func isTagEqualDigest(r harbor.Resource) bool {
	return r.Tag != "" && r.Tag == r.Digest
}
