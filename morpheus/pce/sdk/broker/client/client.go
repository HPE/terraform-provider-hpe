// (C) Copyright 2021-2026 Hewlett Packard Enterprise Development LP

// Package client provides a minimal PCE broker API client. Only the exchange
// that returns Morpheus connection details is implemented; all Morpheus
// resource management is handled by the separate "morpheus" child provider.
package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/models"
)

// SetScmClientToken injects an auth token onto the request context.
type SetScmClientToken func(ctx *context.Context, meta interface{})

// APIClientHandler is the transport contract required by api.do.
type APIClientHandler interface {
	prepareRequest(
		ctx context.Context,
		path string, method string,
		postBody interface{},
		headerParams map[string]string,
		queryParams url.Values,
		formParams url.Values,
		fileName string,
		fileBytes []byte) (localVarRequest *http.Request, err error)
	callAPI(request *http.Request) (*http.Response, error)
	getVersion() int
	getHost() string
}

// APIClient manages communication with the PCE broker API.
type APIClient struct {
	cfg        *Configuration
	cmpVersion int
	meta       interface{}
	tokenFunc  SetScmClientToken
}

// defaultTokenFunc is used until a token function is supplied; it does not
// fetch any token or update the context.
func defaultTokenFunc(_ *context.Context, _ interface{}) {}

// NewAPIClient creates a new API Client, defaulting the HTTP client if unset.
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	return &APIClient{
		cfg:       cfg,
		tokenFunc: defaultTokenFunc,
	}
}

func (c *APIClient) getHost() string {
	return c.cfg.Host
}

// GetCMPDetails exchanges the configured credentials for Morpheus connection
// details. The receiver must be a client whose Host is the PCE broker.
func (c *APIClient) GetCMPDetails(ctx context.Context) (models.TFMorpheusDetails, error) {
	cmpBroker := BrokerAPIService{
		Client: c,
		Cfg:    *c.cfg,
	}

	return cmpBroker.GetMorpheusDetails(ctx)
}

// SetMetaFnAndVersion sets the client token function and the CMP version.
func (c *APIClient) SetMetaFnAndVersion(meta interface{}, version int, fn SetScmClientToken) {
	c.meta = meta
	c.tokenFunc = fn
	c.cmpVersion = version
}

// callAPI performs the request.
//
// The request URL is built from the configured broker host, which is provider
// configuration rather than untrusted input: a Disconnected deployment has to
// be able to point the client at its own on-premise broker, so the host is
// necessarily operator-supplied. gosec reports this as SSRF via taint analysis.
func (c *APIClient) callAPI(request *http.Request) (*http.Response, error) {
	return c.cfg.HTTPClient.Do(request) //nolint:gosec // G704: operator-configured broker host
}

func (c *APIClient) getVersion() int {
	return c.cmpVersion
}

// prepareRequest build the request
// nolint
func (c *APIClient) prepareRequest(
	ctx context.Context,
	path string, method string,
	postBody interface{},
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	fileName string,
	fileBytes []byte) (localVarRequest *http.Request, err error) {
	var body *bytes.Buffer
	// Detect postBody type and post.
	if postBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}
		body, err = setBody(postBody, contentType)
		if err != nil {
			return nil, err
		}
	}
	// add form parameters and file if available.
	if strings.HasPrefix(headerParams["Content-Type"], "multipart/form-data") && len(formParams) > 0 ||
		(len(fileBytes) > 0 && fileName != "") {
		if body != nil {
			return nil, errors.New("cannot specify postBody and multipart form at the same time")
		}
		body = &bytes.Buffer{}
		w := multipart.NewWriter(body)
		for k, v := range formParams {
			for _, iv := range v {
				if strings.HasPrefix(k, "@") { // file
					err = addFile(w, k[1:], iv)
					if err != nil {
						return nil, err
					}
				} else { // form value
					err := w.WriteField(k, iv)
					if err != nil {
						continue
					}
				}
			}
		}
		if len(fileBytes) > 0 && fileName != "" {
			w.Boundary()
			// _, fileNm := filepath.Split(fileName)
			part, err := w.CreateFormFile("file", filepath.Base(fileName))
			if err != nil {
				return nil, err
			}
			_, err = part.Write(fileBytes)
			if err != nil {
				return nil, err
			}
			// Set the Boundary in the Content-Type
			headerParams["Content-Type"] = w.FormDataContentType()
		}
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
		w.Close()
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") &&
		len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("cannot specify postBody and x-www-form-urlencoded form at the same time")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Adding Query Param
	query := url.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	for k, v := range c.cfg.DefaultQueryParams {
		query.Add(k, v)
	}

	// Encode the parameters.
	url.RawQuery = query.Encode()

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	// Add the Authentication Token to header.
	if ctx != nil {
		// Set auth token. This implementation is temporary fix. Need to put
		// this in a goroutine or implement cache in cmp-api
		c.tokenFunc(&ctx, c.meta)
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)
		// Basic HTTP Authentication
		if auth, ok := ctx.Value(ContextBasicAuth).(BasicAuth); ok {
			localVarRequest.SetBasicAuth(auth.UserName, auth.Password)
		}
		// AccessToken Authentication
		if auth, ok := ctx.Value(ContextAccessToken).(string); ok {
			localVarRequest.Header.Set("Authorization", "Bearer "+auth)
		}
	}

	for header, value := range c.cfg.DefaultHeader {
		if value != "" && value != " " {
			localVarRequest.Header.Add(header, value)
		}
	}

	return localVarRequest, nil
}
