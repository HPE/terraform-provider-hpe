// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/utils/httptrace"
)

// RequestRecorder records the API requests the provider sends during a test,
// so a check can assert what an apply did send -- in particular whether a
// write-only value was included -- independently of the appliance. Install it
// with RecordRequests; it is removed when the test ends.
type RequestRecorder struct {
	mu   sync.Mutex
	reqs []RecordedRequest
}

// RecordedRequest is one API request: method, URL path and JSON body (nil when
// there was none or it was not JSON).
type RecordedRequest struct {
	Method string
	Path   string
	Body   map[string]any
}

// RecordRequests installs a recorder for the duration of t. Tests that use it
// must not run in parallel with each other: the observer is process-wide.
func RecordRequests(t *testing.T) *RequestRecorder {
	t.Helper()

	r := &RequestRecorder{}

	httptrace.SetObserver(func(method, path string, body []byte) {
		rec := RecordedRequest{Method: method, Path: path}

		if len(body) > 0 {
			var m map[string]any
			if json.Unmarshal(body, &m) == nil {
				rec.Body = m
			}
		}

		r.mu.Lock()
		r.reqs = append(r.reqs, rec)
		r.mu.Unlock()
	})

	t.Cleanup(func() { httptrace.SetObserver(nil) })

	return r
}

// Reset forgets the requests recorded so far. Call it at the start of a step
// (in a PreConfig) so the step's checks see only that step's requests.
func (r *RequestRecorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.reqs = nil
}

// Writes returns the recorded POST/PUT/PATCH requests whose path contains
// pathPart.
func (r *RequestRecorder) Writes(pathPart string) []RecordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []RecordedRequest

	for _, req := range r.reqs {
		if (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") && strings.Contains(req.Path, pathPart) {
			out = append(out, req)
		}
	}

	return out
}

// hasKey reports whether the dotted key (e.g. "user.password") is present in
// the body, at any nesting under a single object envelope.
func hasKey(body map[string]any, dotted string) bool {
	cur := any(body)

	for _, part := range strings.Split(dotted, ".") {
		m, ok := cur.(map[string]any)
		if !ok {
			return false
		}

		cur, ok = m[part]
		if !ok {
			return false
		}
	}

	return true
}

// ExpectSent is a TestCheckFunc asserting that, since the last Reset, a write
// to a path containing pathPart carried the dotted key (want=true) or that no
// such write did (want=false). The recorded requests are reset afterwards.
//
// Because of that reset, a step that must assert several keys of the same
// write(s) -- one sent, another not -- uses ExpectSentKeys: a second
// ExpectSent in the same step would see no requests at all.
func (r *RequestRecorder) ExpectSent(pathPart, dottedKey string, want bool) resource.TestCheckFunc {
	return r.ExpectSentKeys(pathPart, map[string]bool{dottedKey: want})
}

// ExpectSentKeys is ExpectSent for several dotted keys at once: every key is
// checked against the writes recorded since the last Reset, and the recorded
// requests are reset once, afterwards.
func (r *RequestRecorder) ExpectSentKeys(pathPart string, want map[string]bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		defer r.Reset()

		writes := r.Writes(pathPart)

		var errs []string

		for dottedKey, wantSent := range want {
			sent := false

			for _, w := range writes {
				// An empty key means "any write to the path counts" -- for
				// values that travel in the URL rather than the body.
				if dottedKey == "" || hasKey(w.Body, dottedKey) {
					sent = true

					break
				}
			}

			if sent != wantSent {
				errs = append(errs, fmt.Sprintf("key %q sent=%v, want %v", dottedKey, sent, wantSent))
			}
		}

		if len(errs) > 0 {
			sort.Strings(errs)

			return fmt.Errorf("%d write(s) to %q: %s", len(writes), pathPart, strings.Join(errs, "; "))
		}

		return nil
	}
}
