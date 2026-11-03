// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package httptrace

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

// RequestObserver receives a copy of every API request the provider sends,
// with its body, when one is installed with SetObserver. It exists for
// in-process acceptance tests -- the test harness runs the provider in the
// test's own process -- to assert what was sent (for example that a write-only
// value was, or was not, included in an update) without an appliance-side
// oracle. Nothing is installed in normal operation, so the cost is one nil
// check per request.
type RequestObserver func(method, path string, body []byte)

var (
	observerMu sync.RWMutex
	observer   RequestObserver
)

// SetObserver installs (or, with nil, removes) the request observer.
func SetObserver(o RequestObserver) {
	observerMu.Lock()
	defer observerMu.Unlock()

	observer = o
}

// ObserveRoundTripper wraps transport so the installed observer sees each
// request. It is installed unconditionally; with no observer it is a
// pass-through.
type ObserveRoundTripper struct {
	Transport http.RoundTripper
}

func (t ObserveRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	observerMu.RLock()
	o := observer
	observerMu.RUnlock()

	if o == nil {
		return t.Transport.RoundTrip(req)
	}

	var body []byte

	if req.Body != nil && req.Body != http.NoBody {
		b, err := io.ReadAll(req.Body)
		if err == nil {
			body = b
			req.Body = io.NopCloser(bytes.NewReader(b))
		}
	}

	o(req.Method, req.URL.Path, body)

	return t.Transport.RoundTrip(req)
}
