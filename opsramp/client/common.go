// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// OpsRampClient handles API communication.
type OpsRampClient struct {
	BaseUrl       string
	TenantId      string
	AccessToken   string
	muAccessToken sync.RWMutex
	Client        *http.Client
	Scope         string
}

// OAuthTokenResponse stores API response
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// NewOpsRampClient creates a new client and retrieves the OAuth token.
func NewOpsRampClient(clientID string, clientSecret string, endpoint string, tenant string) (*OpsRampClient, error) {
	rateLimitedTransport := NewRateLimitedTransport(&http.Transport{})

	client := &http.Client{
		Transport: rateLimitedTransport,
	}

	baseUrl := fmt.Sprintf("https://%s", endpoint)
	resource := "/tenancy/auth/oauth/token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	u, _ := url.ParseRequestURI(baseUrl)

	u.Path = resource
	urlStr := u.String()

	resp, err := client.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to retrieve OAuth token")
	}

	var tokenResponse OAuthTokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	opsClient := &OpsRampClient{
		BaseUrl:     baseUrl,
		TenantId:    tenant,
		AccessToken: tokenResponse.AccessToken,
		Client:      client,
		Scope:       "",
	}

	// Detect scope by checking endpoint availability
	_, err = opsClient.GetTenantInfo(tenant)
	if err != nil {
		// If tenant info retrieval fails, default to MSP scope
		opsClient.Scope = "CLIENT"
	} else {
		opsClient.Scope = "MSP"
	}

	return opsClient, nil
}

// Thread-safe getter for AccessToken
func (c *OpsRampClient) GetAccessToken() string {
	c.muAccessToken.RLock()
	defer c.muAccessToken.RUnlock()
	return c.AccessToken
}

// Thread-safe setter for AccessToken
func (c *OpsRampClient) SetAccessToken(token string) {
	c.muAccessToken.Lock()
	defer c.muAccessToken.Unlock()
	c.AccessToken = token
}

func (c *OpsRampClient) NewJsonRequest(method string, apiUrl string, payload []byte) ([]byte, error) {
	validStatuses := map[int]struct{}{
		http.StatusOK:       {},
		http.StatusAccepted: {},
		http.StatusCreated:  {},
	}

	const maxRetries = 3

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Build a fresh request for each attempt.
		payloadReader := strings.NewReader(string(payload))
		req, err := http.NewRequest(method, apiUrl, payloadReader)
		if err != nil {
			return nil, err
		}

		// Add the proper Headers
		req.Header.Set("Authorization", "Bearer "+c.GetAccessToken())
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// Make the Request
		res, err := c.Client.Do(req)
		if err != nil {
			if attempt < maxRetries && shouldRetryNetworkError(err) {
				delay := time.Duration(attempt+1) * time.Second
				time.Sleep(delay)
				continue
			}
			return nil, err
		}

		body, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}

		if _, ok := validStatuses[res.StatusCode]; ok {
			return body, nil
		}

		// Transport-focused retries only. Domain-level retry/reconciliation
		// (e.g. create-then-500 ambiguous outcomes) must be handled by callers.
		if attempt < maxRetries && shouldRetryStatusCode(res.StatusCode) {
			delay := time.Duration(attempt+1) * time.Second
			time.Sleep(delay)
			continue
		}

		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return nil, fmt.Errorf("request failed after retries")
}

func shouldRetryNetworkError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, "timeout") || strings.Contains(errText, "tempor")
}

func shouldRetryStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusInternalServerError,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
