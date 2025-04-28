// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	morpheus "github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	URL              string
	credentials      *credentials
	accessToken      *accessToken
	defaultTransport http.RoundTripper
	ctx              context.Context

	Morpheus *morpheus.APIClient
}
type credentials struct {
	Username string
	Password string
}

type accessToken struct {
	AccessToken string
}

type Config struct {
	URL         string
	InsecureTLS bool
}

func New(ctx context.Context, cfg Config) *Client {
	var client Client

	morpheusCfg := morpheus.NewConfiguration()
	morpheusCfg.Servers[0].URL = cfg.URL
	client.URL = cfg.URL
	client.ctx = ctx
	client.defaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureTLS}, //nolint:gosec
	}

	httpClient := http.DefaultClient
	httpClient.Transport = &client
	morpheusCfg.HTTPClient = httpClient

	client.Morpheus = morpheus.NewAPIClient(morpheusCfg)

	return &client
}

func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	// Send request with no modifications
	resp, err := c.defaultTransport.RoundTrip(req)
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
	tflog.Info(c.ctx, "current token is not authorized, attempt to acquire a new token")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if c.credentials != nil {
		if err := c.GetTokenWithCredentials(ctx); err != nil {
			tflog.Error(c.ctx,
				fmt.Sprintf("could not acquire a new token with user credentials: %s", err))

			return nil, err
		}
	}

	tflog.Info(c.ctx, "new token successfully acquired")

	// Repeat the previous request with the new token
	req.Header.Set("Authorization", c.Morpheus.GetConfig().DefaultHeader["Authorization"])

	return c.defaultTransport.RoundTrip(req)
}

func (c *Client) GetTokenWithCredentials(ctx context.Context) error {
	req := c.Morpheus.AuthenticationAPI.GetAccessToken(ctx)
	req = req.Username(c.credentials.Username)
	req = req.Password(c.credentials.Password)
	req = req.ClientId("morpheus-terraform")
	req = req.GrantType("password")
	req = req.Scope("write")

	token, _, err := c.Morpheus.AuthenticationAPI.GetAccessTokenExecute(req)
	if err != nil {
		msg := `could not authenticate with the Morpheus API using the username:'%s',` +
			`verify that the credentials are correct: %s"`

		return fmt.Errorf(msg, c.credentials.Username, err)
	}

	c.Morpheus.GetConfig().AddDefaultHeader("Authorization", "Bearer "+token.GetAccessToken())

	return nil
}

func (c *Client) SetAccessToken(_ context.Context, token string) error {
	c.accessToken = &accessToken{
		AccessToken: token,
	}

	c.Morpheus.GetConfig().AddDefaultHeader("Authorization", "Bearer "+c.accessToken.AccessToken)

	return nil
}

func (c *Client) SetCredentials(ctx context.Context, username, password string) error {
	c.credentials = &credentials{
		Username: username,
		Password: password,
	}

	return c.GetTokenWithCredentials(ctx)
}
