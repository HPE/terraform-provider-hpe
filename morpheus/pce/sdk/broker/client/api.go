// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	consts "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/common"
)

// all the required validation should be provided as validationFunc
type validationFunc func() error

// parse response as proper model
type jsonPareserFunc func(body []byte) error

type api struct {
	method            string
	path              string
	client            APIClientHandler
	jsonParser        jsonPareserFunc
	validations       []validationFunc
	compatibleVersion string
	// removeVmaasCMPBasePath is used to remove the base path of the vmaas-cmp API, for use by the broker API
	removeVmaasCMPBasePath bool
}

// do will call the API provided. this function will not return any response, but
// response should be catched from jsonParser function itself
func (a *api) do(ctx context.Context, request interface{}, queryParams map[string]string) error {
	// Checked before anything else uses them: a.client is dereferenced below,
	// so an unconfigured api has to be rejected first. These are programmer
	// errors rather than runtime conditions, but a provider SDK must not panic
	// on them: a panic takes down the whole plugin process.
	if a.path == "" || a.method == "" || a.client == nil || a.jsonParser == nil {
		return errors.New("api not properly configured")
	}

	currentVersion, err := parseVersion(a.compatibleVersion)
	if err != nil {
		return fmt.Errorf("failed to parse the compatible version %q: %w", a.compatibleVersion, err)
	}

	if a.client.getVersion() < currentVersion {
		if a.client.getVersion() == 0 {
			return fmt.Errorf("failed to get meta data for cmp-sdk")
		}

		return errVersion
	}
	var (
		localVarHTTPMethod = strings.ToUpper(a.method)
		localVarFileName   string
		localVarFileBytes  []byte
	)

	// Set the path
	if !a.removeVmaasCMPBasePath {
		// Add the base path of the vmaas-cmp API if we are calling the vmaas-cmp API
		a.path = fmt.Sprintf("%s/%s/%s", a.client.getHost(), consts.VmaasCmpAPIBasePath, a.path)
	} else {
		// Don't use the base path of the vmaas-cmp API if we are calling the broker API
		a.path = fmt.Sprintf("%s/%s", a.client.getHost(), a.path)
	}

	for _, validations := range a.validations {
		err := validations()
		if err != nil {
			return err
		}
	}
	// create path and map variables
	localVarPath := a.path

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := getURLValues(queryParams)
	localVarFormParams := url.Values{}

	// set Accept header
	localVarHeaderParams["Accept"] = consts.ContentType

	req, err := a.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, request,
		localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return err
	}

	// Registered before the status check below: that path returns early, and
	// ParseError reads the body without closing it, so deferring any later
	// would leak the body on every error response.
	defer localVarHTTPResponse.Body.Close()

	if localVarHTTPResponse.StatusCode >= 300 {
		return ParseError(localVarHTTPResponse)
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	if err != nil {
		return err
	}

	return a.jsonParser(localVarBody)
}
