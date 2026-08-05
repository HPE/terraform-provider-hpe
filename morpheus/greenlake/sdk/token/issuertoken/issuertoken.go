// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package issuertoken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/HPE/terraform-provider-hpe/morpheus/greenlake/sdk/token/iamversion"
	tokenutil "github.com/HPE/terraform-provider-hpe/morpheus/greenlake/sdk/token/token-util"
)

const (
	retryLimit = 3
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func GenerateToken(
	ctx context.Context,
	clientID,
	clientSecret string,
	identityServiceURL string,
	httpClient tokenutil.HttpClient,
	iamVersion string,
) (string, error) {
	// Generate the parameters and URL for the request
	params, clientURL, err := generateParamsAndURL(clientID, clientSecret, identityServiceURL, iamVersion)
	if err != nil {
		return "", err
	}

	// Create a slice of cancel functions to be returned by the retries
	cancelFuncs := make([]context.CancelFunc, 0)

	// Execute the request, with retries
	resp, err := tokenutil.DoRetries(
		ctx,
		&cancelFuncs,
		func(reqCtx context.Context) (*http.Request, *http.Response, error) {
			// Create the request
			req, errReq := createRequest(reqCtx, params, clientURL)
			if errReq != nil {
				return nil, nil, errReq
			}
			// Close the request after use, i.e. don't reuse the TCP connection
			req.Close = true

			// Execute the request
			respFromDo, errResp := httpClient.Do(req)

			return req, respFromDo, errResp
		},
		retryLimit,
	)
	// Defer execution of cancel functions
	defer executeCancelFuncs(&cancelFuncs)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = tokenutil.ManageHTTPErrorCodes(resp, clientID)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var token TokenResponse

	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// executeCancelFuncs executes all cancel functions in the slice
func executeCancelFuncs(cancelFuncs *[]context.CancelFunc) {
	for _, cancel := range *cancelFuncs {
		cancel()
	}
}

// createRequest creates a new http request
func createRequest(ctx context.Context, params url.Values, clientURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, clientURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// generateParamsAndURL generates the parameters and URL for the request
func generateParamsAndURL(clientID, clientSecret, identityServiceURL, iamVersion string) (url.Values, string, error) {
	params := url.Values{}

	// Add common parameters for an API Client
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	params.Add("grant_type", "client_credentials")

	// Add specific parameters and generate URL for the IAM version
	var clientURL string
	switch iamversion.Version(iamVersion) {
	case iamversion.GLCS:
		params.Add("scope", "hpe-tenant")
		clientURL = fmt.Sprintf("%s/v1/token", identityServiceURL)

	case iamversion.GLP:
		clientURL = identityServiceURL

	default:
		return nil, "", fmt.Errorf("invalid IAM version")
	}

	return params, clientURL, nil
}
