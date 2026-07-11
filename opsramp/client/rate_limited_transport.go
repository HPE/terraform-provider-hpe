// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// bucketState tracks the rate limit state for a single API bucket.
type bucketState struct {
	limit     int
	remaining int
	resetUnix int64
}

// RateLimitedTransport is an http.RoundTripper that respects per-bucket
// rate limits communicated via API response headers.
type RateLimitedTransport struct {
	Base    http.RoundTripper
	mu      sync.Mutex
	buckets map[string]*bucketState
}

// NewRateLimitedTransport creates a transport that tracks per-bucket rate limits.
func NewRateLimitedTransport(base http.RoundTripper) *RateLimitedTransport {
	return &RateLimitedTransport{
		Base:    base,
		buckets: make(map[string]*bucketState),
	}
}

func (t *RateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.waitIfRateLimited()

	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.updateFromHeaders(resp.Header)
	return resp, nil
}

// waitIfRateLimited blocks until all exhausted buckets have reset.
func (t *RateLimitedTransport) waitIfRateLimited() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().Unix()
	for bucket, state := range t.buckets {
		if now >= state.resetUnix {
			delete(t.buckets, bucket)
			continue
		}
		if state.remaining <= 0 {
			waitDuration := time.Until(time.Unix(state.resetUnix, 0))
			if waitDuration > 0 {
				log.Printf("[RateLimit] Bucket %q exhausted (limit %d), waiting %s until reset",
					bucket, state.limit, waitDuration.Round(time.Second))
				t.mu.Unlock()
				time.Sleep(waitDuration)
				t.mu.Lock()
			}
		}
	}
}

// updateFromHeaders parses rate limit headers from the API response and
// stores the bucket state.
func (t *RateLimitedTransport) updateFromHeaders(h http.Header) {
	bucket := cleanHeaderValue(h.Get("X-Ratelimit-Bucket"))
	if bucket == "" {
		return
	}

	limitStr := cleanHeaderValue(h.Get("X-Ratelimit-Limit"))
	remainingStr := cleanHeaderValue(h.Get("X-Ratelimit-Remaining"))
	resetStr := cleanHeaderValue(h.Get("X-Ratelimit-Reset"))

	limit, _ := strconv.Atoi(limitStr)
	remaining, _ := strconv.Atoi(remainingStr)
	resetUnix, _ := strconv.ParseInt(resetStr, 10, 64)

	t.mu.Lock()
	defer t.mu.Unlock()

	t.buckets[bucket] = &bucketState{
		limit:     limit,
		remaining: remaining,
		resetUnix: resetUnix,
	}
}

// cleanHeaderValue strips surrounding quotes and whitespace from header values.
func cleanHeaderValue(s string) string {
	return strings.Trim(s, "\" ")
}
