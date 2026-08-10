// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

//nolint:structcheck
package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/identitytoken"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/issuertoken"
)

// badRequestTokenBody is the zero-valued token response that the issuer returns
// with a 400, and which the resulting error message embeds verbatim. A raw
// string literal keeps it readable, as the JSON is quote-heavy.
const badRequestTokenBody = `{"token_type":"","access_token":"","refresh_token":"",` +
	`"expiry":"0001-01-01T00:00:00Z","expires_in":0,"scope":"","accessTokenOnly":false}`

type testCaseIssuer struct {
	name       string
	ctx        context.Context
	url        string
	statusCode int
	token      issuertoken.TokenResponse
	err        error
}

type testCaseIdentity struct {
	name       string
	ctx        context.Context
	url        string
	statusCode int
	token      identitytoken.TokenResponse
	err        error
}

type testHTTPClient struct {
	statusCode int
	body       interface{}
}

type bodyReadCloser struct {
	body      []byte
	readIndex int
}

func (b *bodyReadCloser) Close() error {
	return nil
}

func (b *bodyReadCloser) Read(p []byte) (n int, err error) {
	if b.readIndex >= len(b.body) {
		err = io.EOF

		return
	}

	n = copy(p, b.body[b.readIndex:])
	b.readIndex += n

	return n, nil
}

func (h *testHTTPClient) Do(req *http.Request) (*http.Response, error) {
	bytes, err := json.Marshal(h.body)
	if err != nil {
		return nil, err
	}

	body := &bodyReadCloser{body: bytes}

	return &http.Response{StatusCode: h.statusCode, Body: body}, nil
}

func createTestClient(
	identityServiceURL,
	passedInToken string,
	statusCode int,
	token interface{},
	vendedServiceClient bool,
) *Client {
	c := New(identityServiceURL, vendedServiceClient, passedInToken)
	if token == nil {
		c.httpClient = &testHTTPClient{
			statusCode: statusCode,
			body:       token,
		}
	} else {
		if vendedServiceClient {
			c.httpClient = &testHTTPClient{
				statusCode: statusCode,
				body:       token.(issuertoken.TokenResponse),
			}
		} else {
			c.httpClient = &testHTTPClient{
				statusCode: statusCode,
				body:       token.(identitytoken.TokenResponse),
			}
		}
	}

	return c
}

func setTestCase() ([]testCaseIssuer, []testCaseIdentity) {
	return []testCaseIssuer{
			{
				name:       "success",
				ctx:        context.Background(),
				url:        "https://hpe-greenlake-tenant.okta.com/oauth2/default",
				statusCode: http.StatusOK,
				token: issuertoken.TokenResponse{
					AccessToken: "access-token",
				},
			},
			{
				name:       "status code 404",
				url:        "https://hpe-greenlake-tenant.okta.com/oauth2/default",
				ctx:        context.Background(),
				statusCode: http.StatusNotFound,
				err:        errors.New("Unexpected status code 404"),
			},
			{
				name:       "status code 400",
				url:        "https://hpe-greenlake-tenant.okta.com/oauth2/default",
				ctx:        context.Background(),
				statusCode: http.StatusBadRequest,
				err:        errors.New("Bad request: {\"token_type\":\"\",\"expires_in\":0,\"access_token\":\"\",\"scope\":\"\"}"),
			},
			{
				name:       "status code 401",
				url:        "https://hpe-greenlake-tenant.okta.com/oauth2/default",
				ctx:        context.Background(),
				statusCode: http.StatusUnauthorized,
				err:        errors.New("Unauthorized access: "),
			},
			{
				name:       "status code 403",
				url:        "https://hpe-greenlake-tenant.okta.com/oauth2/default",
				ctx:        context.Background(),
				statusCode: http.StatusForbidden,
				err:        errors.New("Forbidden: "),
			},
		}, []testCaseIdentity{
			{
				name:       "success",
				ctx:        context.Background(),
				url:        "https://client.greenlake.hpe.com/api/iam/identity",
				statusCode: http.StatusOK,
				token: identitytoken.TokenResponse{
					AccessToken: "access-token",
				},
			},
			{
				name:       "status code 404",
				url:        "https://client.greenlake.hpe.com/api/iam/identity",
				ctx:        context.Background(),
				statusCode: http.StatusNotFound,
				err:        errors.New("Unexpected status code 404"),
			},
			{
				name:       "status code 400",
				url:        "https://client.greenlake.hpe.com/api/iam/identity",
				ctx:        context.Background(),
				statusCode: http.StatusBadRequest,
				err:        errors.New("Bad request: " + badRequestTokenBody),
			},
			{
				name:       "status code 401",
				url:        "https://client.greenlake.hpe.com/api/iam/identity",
				ctx:        context.Background(),
				statusCode: http.StatusUnauthorized,
				err:        errors.New("Unauthorized access: "),
			},
			{
				name:       "status code 403",
				url:        "https://client.greenlake.hpe.com/api/iam/identity",
				ctx:        context.Background(),
				statusCode: http.StatusForbidden,
				err:        errors.New("Forbidden: "),
			},
		}
}

func TestGenerateToken(t *testing.T) {
	t.Parallel()
	var c *Client
	issuertokentc, identitytokentc := setTestCase()

	// Tests for issuertoken package
	for _, testcase := range issuertokentc {
		tc := testcase

		c = createTestClient(tc.url, "", tc.statusCode, tc.token, true)

		token, err := c.GenerateToken(tc.ctx, "", "", "", string(iamversion.GLCS))
		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		}

		assert.Equal(t, tc.token.AccessToken, token)
	}

	// Tests for identitytoken package
	for _, testcase := range identitytokentc {
		tc := testcase

		c = createTestClient(tc.url, "", tc.statusCode, tc.token, false)

		token, err := c.GenerateToken(tc.ctx, "", "", "", string(iamversion.GLCS))
		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		}

		assert.Equal(t, tc.token.AccessToken, token)
	}
}

func TestGenerateTokenPassedInToken(t *testing.T) {
	t.Parallel()
	c := createTestClient("", "testToken", http.StatusAccepted, nil, true)

	token, err := c.GenerateToken(context.Background(), "", "", "", "")
	assert.Equal(t, "testToken", token)
	assert.NoError(t, err)
}
