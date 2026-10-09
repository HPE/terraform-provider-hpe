// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

type CredsRoundTripper struct {
	baseTransport http.RoundTripper
	// tokenMu guards the cached token (the SDK configuration's default
	// Authorization header): every read of it, and every re-authentication.
	tokenMu  sync.Mutex
	client   *sdk.APIClient
	username string
	password string
}

// GetToken authenticates with the username and password and caches the
// access token as the SDK configuration's default Authorization header, which
// the SDK then adds to every request it makes.
//
// It must be called with tokenMu held.
//
// WHY the header is removed first: the SDK adds the default headers to the
// token request too. When GetToken runs to replace an expired token, that
// would send the old bearer token along with the password grant -- and the
// Morpheus token endpoint refuses any request that carries a bearer token,
// valid or expired: its authorization-server filter chain tries to decode the
// bearer as a JWT, Morpheus access tokens are opaque, so the request is
// answered 401 ("invalid_token ... Malformed token") before the grant is even
// read (seen on 9.0.0 and 9.1.0.2). The refresh could therefore never
// succeed, and every operation longer than the token's lifetime (1439 s for
// this client, about 24 minutes) failed with "could not authenticate".
// Clearing the header before the call also means a failed refresh leaves no
// stale token behind for the next request to retry with.
func (c *CredsRoundTripper) GetToken(ctx context.Context) error {
	delete(c.client.GetConfig().DefaultHeader, "Authorization")

	req := c.client.AuthenticationAPI.GetAccessToken(ctx)
	req = req.Username(c.username)
	req = req.Password(c.password)
	req = req.ClientId("morpheus-terraform")
	req = req.GrantType("password")
	req = req.Scope("write")

	token, _, err := c.client.AuthenticationAPI.GetAccessTokenExecute(req)
	if err != nil {
		msg := `could not authenticate with the Morpheus API using ` +
			`the username:'%s',verify that the credentials are correct: %s"`

		return fmt.Errorf(msg, c.username, err)
	}

	if token == nil || token.AccessToken == nil {
		return fmt.Errorf(
			"could not authenticate with the Morpheus API using the username:'%s': response missing access token", c.username)
	}

	c.client.GetConfig().AddDefaultHeader(
		"Authorization", "Bearer "+*token.AccessToken,
	)

	return nil
}

func (c *CredsRoundTripper) cachedAuthHeader() string {
	return c.client.GetConfig().DefaultHeader["Authorization"]
}

// InitReqAuthHeader performs first time initialization
func (c *CredsRoundTripper) InitReqAuthHeader(req *http.Request) error {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if req.Header.Get("Authorization") != "" {
		return nil
	}

	hdr := c.cachedAuthHeader()
	if hdr != "" {
		req.Header.Set("Authorization", hdr)

		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := c.GetToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.cachedAuthHeader())

	return nil
}

// refreshToken replaces the cached token after a request was refused with
// 401, and returns the new Authorization header value.
//
// Concurrent requests that were all refused with the same expired token must
// not each log in again: whoever takes the lock first re-authenticates, and
// the others find the token already replaced and reuse it.
func (c *CredsRoundTripper) refreshToken(ctx context.Context, rejected string) (string, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if hdr := c.cachedAuthHeader(); hdr != "" && hdr != rejected {
		return hdr, nil
	}

	tokenCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tflog.Debug(ctx, "refreshing token")

	if err := c.GetToken(tokenCtx); err != nil {
		return "", err
	}

	tflog.Debug(ctx, "new token successfully acquired")

	return c.cachedAuthHeader(), nil
}

func (c *CredsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	err := c.InitReqAuthHeader(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.baseTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Request authorized, keep going
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	// The refused response is not returned to anyone; release its connection.
	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	// The token was refused (expired, most likely): log in again and repeat
	// the request once with the new token. Any further retrying belongs to
	// the callers, not here.
	hdr, err := c.refreshToken(req.Context(), req.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", hdr)

	// The first attempt consumed the body; the retry needs it again. The SDK
	// builds its requests from an in-memory buffer, so GetBody is set.
	if req.GetBody != nil && req.Body != nil && req.Body != http.NoBody {
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("retrying request after re-authentication: %w", err)
		}

		req.Body = body
	}

	return c.baseTransport.RoundTrip(req)
}

func NewCredsRoundTripper(
	_ context.Context,
	transport http.RoundTripper,
	url string,
	username string,
	password string,
) http.RoundTripper {
	url = NormalizeBaseURL(url)
	morpheusCfg := sdk.NewConfiguration()
	morpheusCfg.Servers[0].URL = url

	client := sdk.NewAPIClient(morpheusCfg)

	morpheusCfg.HTTPClient = &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}

	rt := CredsRoundTripper{
		baseTransport: transport,
		client:        client,
		username:      username,
		password:      password,
	}

	return &rt
}
