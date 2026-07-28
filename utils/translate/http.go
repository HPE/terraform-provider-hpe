// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// doCreate handles POST operations.
func (c *Client) doCreate(ctx context.Context, req Request, cc *CompiledConfig) (*Response, error) {
	if req.Model == nil {
		return nil, fmt.Errorf("create requires a model")
	}

	body, err := Marshal(ctx, req.Model, cc)
	if err != nil {
		return nil, fmt.Errorf("create marshal: %w", err)
	}

	path := c.resolvePath(cc, "create", req)

	return c.doRequest(ctx, http.MethodPost, path, body, cc, req.QueryParams)
}

// doRead handles GET operations.
func (c *Client) doRead(ctx context.Context, req Request, cc *CompiledConfig) (*Response, error) {
	path := c.resolvePath(cc, "read", req)

	return c.doRequest(ctx, http.MethodGet, path, nil, cc, req.QueryParams)
}

// doUpdate handles PUT operations.
func (c *Client) doUpdate(ctx context.Context, req Request, cc *CompiledConfig) (*Response, error) {
	if req.Model == nil {
		return nil, fmt.Errorf("update requires a model")
	}

	body, err := Marshal(ctx, req.Model, cc)
	if err != nil {
		return nil, fmt.Errorf("update marshal: %w", err)
	}

	path := c.resolvePath(cc, "update", req)

	return c.doRequest(ctx, http.MethodPut, path, body, cc, req.QueryParams)
}

// doDelete handles DELETE operations.
func (c *Client) doDelete(ctx context.Context, req Request, cc *CompiledConfig) (*Response, error) {
	path := c.resolvePath(cc, "delete", req)

	return c.doRequest(ctx, http.MethodDelete, path, nil, cc, req.QueryParams)
}

// doList handles GET list operations.
func (c *Client) doList(ctx context.Context, req Request, cc *CompiledConfig) (*Response, error) {
	path := c.resolvePath(cc, "list", req)

	return c.doRequest(ctx, http.MethodGet, path, nil, cc, nil)
}

// doRequest performs the actual HTTP request.
func (c *Client) doRequest(
	ctx context.Context, method, path string, body map[string]any,
	cc *CompiledConfig, queryParams map[string]string,
) (*Response, error) {
	reqURL := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}

		bodyReader = bytes.NewReader(jsonBytes)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Append query parameters
	if len(queryParams) > 0 {
		q := httpReq.URL.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}

		httpReq.URL.RawQuery = q.Encode()
	}

	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	httpReq.Header.Set("Accept", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		cfg:        cc,
	}
	// Parse response body if present
	if len(respBody) > 0 {
		var raw map[string]any
		if err := json.Unmarshal(respBody, &raw); err != nil {
			// If it's not JSON, that's OK for DELETE responses
			if method != http.MethodDelete {
				return nil, fmt.Errorf("parsing response (status %d): %w\nbody: %s",
					httpResp.StatusCode, err, string(respBody))
			}
		}

		resp.raw = raw
	}

	// Check for HTTP errors
	if httpResp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("API returned status %d", httpResp.StatusCode)
		if resp.raw != nil {
			if msg, ok := resp.raw["msg"].(string); ok {
				errMsg += ": " + msg
			} else if msg, ok := resp.raw["message"].(string); ok {
				errMsg += ": " + msg
			}

			if errors, ok := resp.raw["errors"]; ok {
				errMsg += fmt.Sprintf(" errors: %v", errors)
			}
		}

		return resp, fmt.Errorf("%s", errMsg)
	}

	return resp, nil
}

// resolvePath determines the API path for a given operation.
func (c *Client) resolvePath(cc *CompiledConfig, operation string, req Request) string {
	// Check if paths are explicitly configured
	if cc.raw.Paths != nil {
		if pathTemplate, ok := cc.raw.Paths[operation]; ok {
			// Parse "METHOD /path" format — we only need the path part
			parts := strings.Fields(pathTemplate)
			path := pathTemplate
			if len(parts) == 2 {
				path = parts[1]
			}

			// Replace path parameters
			path = strings.ReplaceAll(path, "{id}", fmt.Sprintf("%d", req.ID))
			path = strings.ReplaceAll(path, "{parentId}", fmt.Sprintf("%d", req.ParentID))

			return path
		}
	}

	// Fallback: construct from envelope key and resource name
	// This is a reasonable default for the Morpheus API pattern
	if cc.raw.Envelope != nil && cc.raw.Envelope.Request != "" {
		basePath := "/api/" + guessAPIPath(cc.raw.Envelope.Request)
		switch operation {
		case "create":
			return basePath
		case "list":
			return basePath
		case "read", "update", "delete":
			return fmt.Sprintf("%s/%d", basePath, req.ID)
		}
	}

	// Last resort: use resource name from the request
	basePath := "/api/" + req.Resource + "s"
	switch operation {
	case "create", "list":
		return basePath
	default:
		return fmt.Sprintf("%s/%d", basePath, req.ID)
	}
}

// guessAPIPath converts an envelope key to an API path segment.
// e.g., "zone" → "zones", "environment" → "environments"
func guessAPIPath(envelope string) string {
	if strings.HasSuffix(envelope, "s") {
		return envelope
	}

	return envelope + "s"
}
