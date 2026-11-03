// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import "testing"

// TestUnitClusterHostHashExtractor covers the cluster hash extractor: a present
// per-host map, an empty-but-present map (hosts read, none holding a secret), a
// not-present carrier (a partial read), a nil carrier and a wrong-typed
// response. The extractor emits canonical JSON, so the expected hash is the
// JSON encoding of the map (keys sorted by encoding/json).
func TestUnitClusterHostHashExtractor(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		resp        any
		wantHash    string
		wantPresent bool
	}{
		{
			"present map",
			&clusterHostDigest{hashes: map[string]string{"1": "abc", "2": "def"}, present: true},
			`{"1":"abc","2":"def"}`, true,
		},
		{
			"hosts read but none hold a secret: present and empty hashes",
			&clusterHostDigest{hashes: map[string]string{"1": "", "2": ""}, present: true},
			`{"1":"","2":""}`, true,
		},
		{
			"partial read: not present",
			&clusterHostDigest{present: false},
			"", false,
		},
		{"nil carrier: not present", (*clusterHostDigest)(nil), "", false},
		{"wrong type: not present", "not a carrier", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hash, present := clusterHostHashExtractor(tc.resp)
			if hash != tc.wantHash || present != tc.wantPresent {
				t.Fatalf("clusterHostHashExtractor() = (%q, %v), want (%q, %v)",
					hash, present, tc.wantHash, tc.wantPresent)
			}
		})
	}
}
