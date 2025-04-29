// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package auth

import (
	"context"
	"crypto/tls"
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
	token string,
) http.RoundTripper {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	}

	rt := TokenRoundTripper{
		baseTransport: t,
		token:         token,
	}

	return &rt
}
