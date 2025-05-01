// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package auth

import (
	"context"
	"net/http"
	"sync"
)

type TokenRoundTripper struct {
	baseTransport http.RoundTripper
	tokenMu       sync.Mutex
	initToken     bool
	token         string
}

func (t *TokenRoundTripper) InitAuthHeader(req *http.Request) error {
	t.tokenMu.Lock()
	defer t.tokenMu.Unlock()

	if t.initToken {
		return nil
	}

	t.initToken = true
	req.Header.Set("Authorization", "Bearer "+t.token)

	return nil
}

func (t *TokenRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	t.InitAuthHeader(req)

	return t.baseTransport.RoundTrip(req)
}

func NewTokenRoundTripper(
	_ context.Context,
	transport http.RoundTripper,
	token string,
) http.RoundTripper {
	rt := TokenRoundTripper{
		baseTransport: transport,
		token:         token,
	}

	return &rt
}
