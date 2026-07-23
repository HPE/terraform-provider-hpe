// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseSearchPage feeds a real global Search API response (captured from an
// appliance) through the sweeper's decoder. The payload contains dateCreated
// values with a non-RFC3339 offset ("...+0000") that break the SDK's typed
// decode; parseSearchPage must ignore them and still extract id/name/type.
func TestParseSearchPage(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("testdata", "search_response.json"))
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	hits, total, err := parseSearchPage(body)
	if err != nil {
		t.Fatalf("parseSearchPage errored on non-RFC3339 dateCreated: %v", err)
	}

	if total != 9 {
		t.Errorf("total = %d, want 9", total)
	}

	if len(hits) != 9 {
		t.Fatalf("len(hits) = %d, want 9", len(hits))
	}

	want := SearchHit{ID: "402", Name: "qatf-ttestacc-haproxy-lb-full", Type: "Instance"}
	if hits[0] != want {
		t.Errorf("hits[0] = %+v, want %+v", hits[0], want)
	}

	var instances int
	for _, h := range hits {
		if h.Type == "Instance" {
			instances++
		}
	}
	if instances != 4 {
		t.Errorf("Instance-type hits = %d, want 4", instances)
	}
}

func TestParseSearchPageInvalidJSON(t *testing.T) {
	if _, _, err := parseSearchPage([]byte(`{not valid`)); err == nil {
		t.Error("expected an error for invalid JSON, got nil")
	}
}
