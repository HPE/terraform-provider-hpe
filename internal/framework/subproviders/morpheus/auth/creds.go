// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CredsRoundTripper struct {
	baseTransport http.RoundTripper
	tokenMu       sync.Mutex
	client        *sdk.APIClient
	username      string
	password      string
}

func (c *CredsRoundTripper) GetToken(ctx context.Context) error {
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

	c.client.GetConfig().AddDefaultHeader(
		"Authorization", "Bearer "+token.GetAccessToken(),
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

	// Original request failed with an Unauthorized status, so we should try to
	// generate a new token. We might want to add some retrying here in the
	// future, however it might not be necessary depending on how the
	// retrying is implemented for the API outside of the client.
	tflog.Debug(req.Context(), "refreshing token")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	tflog.Debug(req.Context(), "new token successfully acquired")

	// Repeat the previous request with the new token
	req.Header.Set("Authorization", c.client.GetConfig().DefaultHeader["Authorization"])

	return c.baseTransport.RoundTrip(req)
}

func NewCredsRoundTripper(
	_ context.Context,
	transport http.RoundTripper,
	url string,
	username string,
	password string,
) http.RoundTripper {
	morpheusCfg := sdk.NewConfiguration()
	morpheusCfg.Servers[0].URL = url

	client := sdk.NewAPIClient(morpheusCfg)

	morpheusCfg.HTTPClient = &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	rt := CredsRoundTripper{
		baseTransport: transport,
		client:        client,
		username:      username,
		password:      password,
	}

	return &rt
}
