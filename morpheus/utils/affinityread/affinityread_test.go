// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package affinityread_test

import (
	"net/http"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinityread"
)

// entry stands in for a generated list element. The real cloud and cluster
// types differ while their shapes do not, which is why the helpers take an
// accessor.
type entry struct {
	id   int64
	null bool
}

func entryID(e entry) (int64, bool) {
	if e.null {
		return 0, false
	}

	return e.id, true
}

func TestIsSingleItemRenderFailure(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		resp *http.Response
		want bool
	}{
		{"500 is the defect", &http.Response{StatusCode: http.StatusInternalServerError}, true},
		{"404 is not", &http.Response{StatusCode: http.StatusNotFound}, false},
		{"200 is not", &http.Response{StatusCode: http.StatusOK}, false},
		{"502 is not, the check is narrow", &http.Response{StatusCode: http.StatusBadGateway}, false},
		{"nil response is not", nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := affinityread.IsSingleItemRenderFailure(tc.resp); got != tc.want {
				t.Errorf("IsSingleItemRenderFailure = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestIDsFromListEmptyIsNotNil pins the distinction the create path depends on:
// a cloud or cluster with no groups yields an empty map, NOT nil. nil is
// reserved for "the listing could not be read", and conflating the two would
// turn a recoverable create into a failed one.
func TestIDsFromListEmptyIsNotNil(t *testing.T) {
	t.Parallel()

	got := affinityread.IDsFromList([]entry{}, entryID)
	if got == nil {
		t.Fatal("IDsFromList returned nil for an empty listing; nil means unreadable")
	}

	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestIDsFromList(t *testing.T) {
	t.Parallel()

	got := affinityread.IDsFromList(
		[]entry{{id: 1}, {id: 2}, {null: true}, {id: 3}}, entryID,
	)

	if len(got) != 3 {
		t.Fatalf("len = %d, want 3 (the nil id must be skipped)", len(got))
	}

	for _, want := range []int64{1, 2, 3} {
		if _, ok := got[want]; !ok {
			t.Errorf("id %d missing", want)
		}
	}
}

// TestNewIDFromList covers the recovery of a created group's id when the create
// response could not be rendered.
//
// The "two new" row is the one that matters: if another client created a group
// at the same moment, the choice would be a guess, and adopting a stranger's
// group into state is worse than failing.
func TestNewIDFromList(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		before []int64
		after  []entry
		wantID int64
		wantOK bool
	}{
		{
			name:   "exactly one new id is the created group",
			before: []int64{1, 2},
			after:  []entry{{id: 1}, {id: 2}, {id: 7}},
			wantID: 7,
			wantOK: true,
		},
		{
			name:   "first group on an empty cluster",
			before: nil,
			after:  []entry{{id: 4}},
			wantID: 4,
			wantOK: true,
		},
		{
			name:   "no new id means nothing was created",
			before: []int64{1, 2},
			after:  []entry{{id: 1}, {id: 2}},
			wantOK: false,
		},
		{
			name:   "two new ids are ambiguous and must not be guessed",
			before: []int64{1},
			after:  []entry{{id: 1}, {id: 8}, {id: 9}},
			wantOK: false,
		},
		{
			name:   "nil ids are skipped, not counted as new",
			before: []int64{1},
			after:  []entry{{id: 1}, {null: true}, {id: 5}},
			wantID: 5,
			wantOK: true,
		},
		{
			name:   "a group removed alongside one added is still unambiguous",
			before: []int64{1, 2},
			after:  []entry{{id: 2}, {id: 6}},
			wantID: 6,
			wantOK: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before := make(map[int64]struct{}, len(tc.before))
			for _, id := range tc.before {
				before[id] = struct{}{}
			}

			gotID, gotOK := affinityread.NewIDFromList(tc.after, before, entryID)
			if gotOK != tc.wantOK {
				t.Fatalf("ok = %v, want %v", gotOK, tc.wantOK)
			}

			if gotOK && gotID != tc.wantID {
				t.Errorf("id = %d, want %d", gotID, tc.wantID)
			}
		})
	}
}

func TestServersFromList(t *testing.T) {
	t.Parallel()

	type ag struct {
		id      int64
		servers []int64
	}

	items := []ag{
		{id: 1, servers: []int64{10, 11}},
		{id: 2, servers: nil},
	}

	agID := func(a ag) (int64, bool) { return a.id, true }
	agServers := func(a ag) []int64 { return a.servers }

	t.Run("membership is returned for the matching group", func(t *testing.T) {
		t.Parallel()

		got, ok := affinityread.ServersFromList(items, 1, agID, agServers)
		if !ok {
			t.Fatal("group 1 not found")
		}

		if len(got) != 2 {
			t.Errorf("len = %d, want 2", len(got))
		}
	})

	// An empty group must come back as an empty slice, not nil: nil would read
	// as "unknown membership" and the caller would preserve stale state.
	t.Run("an empty group yields an empty slice not nil", func(t *testing.T) {
		t.Parallel()

		got, ok := affinityread.ServersFromList(items, 2, agID, agServers)
		if !ok {
			t.Fatal("group 2 not found")
		}

		if got == nil {
			t.Error("got nil, want an empty slice")
		}

		if len(got) != 0 {
			t.Errorf("len = %d, want 0", len(got))
		}
	})

	t.Run("an absent group is reported missing", func(t *testing.T) {
		t.Parallel()

		if _, ok := affinityread.ServersFromList(items, 99, agID, agServers); ok {
			t.Error("group 99 reported found")
		}
	})
}
