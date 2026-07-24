// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package tokenutil

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/errors"
	jose "github.com/go-jose/go-jose/v3"
)

// Token a jwt token format
type Token struct {
	Issuer           string `json:"iss"`
	Subject          string `json:"sub"`
	Expiry           int64  `json:"exp"`
	IssuedAt         int64  `json:"iat"`
	Type             string `json:"typ"`
	Nonce            string `json:"nonce"`
	AtHash           string `json:"at_hash"`
	ClientID         string `json:"cid,omitempty"`
	UserID           string `json:"uid,omitempty"`
	TenantID         string `json:"tenantId"`
	AuthorizedParty  string `json:"azp"`
	KeycloakClientID string `json:"clientId"`
	IsHPE            bool   `json:"isHPE"`
}

//nolint:stylecheck,golint,revive
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DecodeAccessToken decodes the accessToken offline
//
//nolint:gocritic
func DecodeAccessToken(rawToken string) (Token, error) {
	_, err := jose.ParseSigned(rawToken)
	if err != nil {
		return Token{}, fmt.Errorf("oidc: malformed jwt: %w", err)
	}

	// Throw out tokens with invalid claims before trying to verify the token. This lets
	// us do cheap checks before possibly re-syncing keys.
	payload, err := parseJWT(rawToken)
	if err != nil {
		log.Fatalf("oidc: malformed jwt: %v", err)

		return Token{}, fmt.Errorf("oidc: malformed jwt: %w", err)
	}
	var token Token
	if err := json.Unmarshal(payload, &token); err != nil {
		log.Fatalf("oidc: failed to unmarshal claims: %v", err)

		return Token{}, fmt.Errorf("oidc: failed to unmarshal claims: %w", err)
	}

	if token.UserID != "" {
		// User token
		token.Subject = "users/" + token.UserID
	} else if token.ClientID != "" || token.KeycloakClientID != "" {
		token.Subject = "clients/" + token.Subject
	} else {
		// TODO This is just so that Keycloak tokens continue to work. Remove after keycloak is gone
		token.Subject = "users/" + token.Subject
	}

	return token, nil
}

func DoRetries(
	ctx context.Context,
	cancelFuncs *[]context.CancelFunc,
	call func(ctx context.Context) (*http.Request, *http.Response, error),
	retries int,
) (*http.Response, error) {
	var req *http.Request
	var resp *http.Response
	var err error

	for {
		// If retries are exhausted, return an error
		if retries == 0 {
			return resp, errors.MakeErrInternalError(errors.ErrorResponse{
				ErrorCode: "ErrGenerateTokenRetryLimitExceeded",
				Message:   "Retry limit exceeded"})
		}

		// Create a new context with a timeout
		ctxWithTimeout, cancel := createContextWithTimeout(ctx)

		// Add the cancel function to the list of cancel functions
		*cancelFuncs = append(*cancelFuncs, cancel)

		// Execute the request
		req, resp, err = call(ctxWithTimeout)

		// If the error is due to a context timeout, retry the request
		if req != nil && req.Context().Err() == context.DeadlineExceeded {
			retries = sleepAndDecrementRetries(retries)

			continue
		}

		// For all other errors, return the error
		if err != nil {
			return resp, err
		}

		// If the status code is not retryable, return the response
		if !isStatusRetryable(resp.StatusCode) {
			return resp, nil
		}

		retries = sleepAndDecrementRetries(retries)
	}
}

func createContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return context.WithTimeout(context.Background(), 3*time.Second)
	}

	return context.WithTimeout(ctx, 3*time.Second)
}

func sleepAndDecrementRetries(retries int) int {
	log.Printf("Retrying request, retries left: %v", retries)
	time.Sleep(5 * time.Second)

	return retries - 1
}

func ManageHTTPErrorCodes(resp *http.Response, clientID string) error {
	var err error

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Bad request: %v", string(body))
		err = errors.MakeErrBadRequest(errors.ErrorResponse{
			ErrorCode: "ErrGenerateTokenBadRequest",
			Message:   msg,
		})

		return err
	case http.StatusForbidden:
		err = errors.MakeErrForbidden(clientID)

		return err
	case http.StatusUnauthorized:
		err = errors.MakeErrUnauthorized(clientID)

		return err
	default:
		msg := fmt.Sprintf("Unexpected status code %v", resp.StatusCode)
		err = errors.MakeErrInternalError(errors.ErrorResponse{
			ErrorCode: "ErrGenerateTokenUnexpectedResponseCode",
			Message:   msg,
		})

		return err
	}
}

func isStatusRetryable(statusCode int) bool {
	if statusCode == http.StatusInternalServerError || statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusBadGateway {
		return true
	}

	return false
}

func parseJWT(p string) ([]byte, error) {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("oidc: malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt payload: %w", err)
	}

	return payload, nil
}
