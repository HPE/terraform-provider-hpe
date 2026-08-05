// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package identitytoken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tokenutil "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/token-util"
)

const (
	retryLimit = 3
)

type GenerateTokenInput struct {
	TenantID     string `json:"tenant_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type TokenResponse struct {
	TokenType       string    `json:"token_type"`
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	Expiry          time.Time `json:"expiry"`
	ExpiresIn       int       `json:"expires_in"`
	Scope           string    `json:"scope"`
	AccessTokenOnly bool      `json:"accessTokenOnly"`
}

func GenerateToken(
	ctx context.Context,
	tenantID,
	clientID,
	clientSecret string,
	identityServiceURL string,
	httpClient tokenutil.HttpClient,
) (string, error) {
	params := GenerateTokenInput{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	url := fmt.Sprintf("%s/v1/token", identityServiceURL)

	b, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// Create a slice of cancel functions to be returned by the retries
	cancelFuncs := make([]context.CancelFunc, 0)

	resp, err := tokenutil.DoRetries(
		ctx,
		&cancelFuncs,
		func(reqCtx context.Context) (*http.Request, *http.Response, error) {
			req, errReq := http.NewRequestWithContext(reqCtx, http.MethodPost, url, strings.NewReader(string(b)))
			if errReq != nil {
				return nil, nil, errReq
			}
			req.Header.Set("Content-Type", "application/json")
			respFromDo, errResp := httpClient.Do(req)

			return req, respFromDo, errResp
		},
		retryLimit,
	)
	// Defer execution of cancelFuncs
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
