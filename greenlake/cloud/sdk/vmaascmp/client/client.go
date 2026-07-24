// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

//go:generate go run github.com/golang/mock/mockgen -source ./client.go -package client -destination ./client_mock.go

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type SetScmClientToken func(ctx *context.Context, meta interface{})

type APIClientHandler interface {
	ChangeBasePath(path string)
	prepareRequest(
		ctx context.Context,
		path string, method string,
		postBody interface{},
		headerParams map[string]string,
		queryParams url.Values,
		formParams url.Values,
		fileName string,
		fileBytes []byte) (localVarRequest *http.Request, err error)
	decode(v interface{}, b []byte, contentType string) (err error)
	callAPI(request *http.Request) (*http.Response, error)
	SetMeta(meta interface{}, fn SetScmClientToken) error
	getVersion() int
	getHost() string
	// The next two methods are for use when creating the Broker API client
	// SetMetaFnAndVersion is used to set the client token function in meta and the SCM version for the Broker client
	SetMetaFnAndVersion(meta interface{}, version int, fn SetScmClientToken)
	// GetSCMVersion returns the SCM version for use when creating the Broker client
	GetSCMVersion() int
	SetHost(host string)
	// GetCMPDetails here the client is the one which has broker host set
	GetCMPDetails(ctx context.Context) (models.TFMorpheusDetails, error)
	SetCMPVersion(ctx context.Context) (err error)
}

// APIClient manages communication with the GreenLake Private Cloud VMaaS CMP API API v1.0.0
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg         *Configuration
	cmpVersion  int
	meta        interface{}
	tokenFunc   SetScmClientToken
	CMPToken    string
	TokenExpiry int64
}

// defaultTokenFunc will use while defining httpClient. defaultTokenFunc
// will not fetch any token or update context.
func defaultTokenFunc(ctx *context.Context, meta interface{}) {}

// NewAPIClient creates a new API Client. Requires a userAgent string describing your application.
// optionally a custom http.Client to allow for advanced features such as caching.
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	c := &APIClient{
		cfg:       cfg,
		tokenFunc: defaultTokenFunc,
	}

	return c
}

func (c *APIClient) getHost() string {
	return c.cfg.Host
}
func (c *APIClient) SetHost(host string) {
	c.cfg.Host = host
}
func (c *APIClient) SetMeta(meta interface{}, fn SetScmClientToken) error {
	c.meta = meta
	c.tokenFunc = fn
	// if cmp version already set then skip
	if c.cmpVersion != 0 {
		return nil
	}
	// initialize cmp status client and get setup/check
	// and set version
	cmpClient := CmpStatus{
		Client: c,
		Cfg:    *c.cfg,
	}
	// Get status of cmp
	statusResp, err := cmpClient.GetCmpVersion(context.Background())
	if err != nil {
		return err
	}
	versionInt, err := parseVersion(statusResp.Appliance.BuildVersion)
	if err != nil {
		return fmt.Errorf("failed to parse cmp build, error: %w", err)
	}
	c.cmpVersion = versionInt

	return nil
}

// GetCMPDetails here APIClient is the brokerClient
func (c *APIClient) GetCMPDetails(ctx context.Context) (models.TFMorpheusDetails, error) {
	cmpBroker := BrokerAPIService{
		Client: c,
		Cfg:    *c.cfg,
	}
	return cmpBroker.GetMorpheusDetails(ctx)

}

func (c *APIClient) SetCMPVersion(ctx context.Context) (err error) {
	if c.cmpVersion != 0 {
		return nil
	}
	cmpClient := CmpStatus{
		Client: c,
		Cfg:    *c.cfg,
	}
	// Get status of cmp
	statusResp, err := cmpClient.GetCmpVersion(ctx)
	if err != nil {
		return
	}
	versionInt, err := parseVersion(statusResp.Appliance.BuildVersion)
	if err != nil {
		return fmt.Errorf("failed to parse cmp build, error: %w", err)
	}
	c.cmpVersion = versionInt
	return
}
func (c *APIClient) SetMetaFnAndVersion(meta interface{}, version int, fn SetScmClientToken) {
	c.meta = meta
	c.tokenFunc = fn
	c.cmpVersion = version
}

// callAPI do the request.
func (c *APIClient) callAPI(request *http.Request) (*http.Response, error) {
	return c.cfg.HTTPClient.Do(request)
}

// Change base path to allow switching to mocks
func (c *APIClient) ChangeBasePath(path string) {
	c.cfg.BasePath = path
}

func (c *APIClient) getVersion() int {
	return c.cmpVersion
}

func (c *APIClient) GetSCMVersion() int {
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

func (c *APIClient) decode(content interface{}, bytes []byte, contentType string) (err error) {
	if strings.Contains(contentType, "application/xml") {
		if err = xml.Unmarshal(bytes, content); err != nil {
			return err
		}

		return nil
	} else if strings.Contains(contentType, "application/json") {
		if err = json.Unmarshal(bytes, content); err != nil {
			return err
		}

		return nil
	}

	return errors.New("undefined response type")
}
