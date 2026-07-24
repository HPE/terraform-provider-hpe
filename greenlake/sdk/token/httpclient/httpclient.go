// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package httpclient

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/identitytoken"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/issuertoken"
	tokenutil "github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/token-util"
)

type Client struct {
	passedInToken       string
	identityServiceURL  string
	httpClient          tokenutil.HttpClient
	vendedServiceClient bool
}

// New creates a new identity Client object
func New(identityServiceURL string, vendedServiceClient bool, passedInToken string) *Client {
	client := &http.Client{Timeout: 120 * time.Second}
	identityServiceURL = strings.TrimRight(identityServiceURL, "/")

	return &Client{
		passedInToken:       passedInToken,
		identityServiceURL:  identityServiceURL,
		httpClient:          client,
		vendedServiceClient: vendedServiceClient,
	}
}

func (c *Client) GenerateToken(ctx context.Context, tenantID, clientID, clientSecret, iamVersion string) (string, error) {
	// we don't have a passed-in token, so we need to actually generate a token
	if c.passedInToken == "" {
		if c.vendedServiceClient {
			token, err := issuertoken.GenerateToken(
				ctx, clientID, clientSecret, c.identityServiceURL, c.httpClient, iamVersion)

			return token, err
		}

		token, err := identitytoken.GenerateToken(ctx, tenantID, clientID, clientSecret, c.identityServiceURL, c.httpClient)

		return token, err
	}

	// we have a passed-in token, return it
	return c.passedInToken, nil
}
