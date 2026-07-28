// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package translate provides a translation layer that converts between
// Terraform resource model structs and the Morpheus API JSON format.
//
// It replaces the generated OpenAPI SDK by using config.yaml transformation
// rules (the same format used by the code-spec pipeline) to reshape data
// at runtime. This eliminates hundreds of lines of manual field mapping.
package translate

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Operation represents a CRUD operation type.
type Operation int

const (
	Create Operation = iota
	Read
	Update
	Delete
	List
)

func (o Operation) String() string {
	switch o {
	case Create:
		return "create"
	case Read:
		return "read"
	case Update:
		return "update"
	case Delete:
		return "delete"
	case List:
		return "list"
	default:
		return "unknown"
	}
}

// Request describes what the translation layer should do.
type Request struct {
	// Operation is the CRUD operation to perform.
	Operation Operation
	// Resource is the resource name (e.g., "cloud", "environment").
	Resource string
	// Model is the Terraform model struct (for Create/Update).
	Model any
	// ID is the resource ID (for Read/Update/Delete).
	ID int64
	// Plan is the Terraform plan model (for Read, enables plan preservation).
	Plan any
	// ParentID is used for sub-resources that require a parent path parameter.
	ParentID int64
	// QueryParams are additional query parameters to append to the URL.
	QueryParams map[string]string
}

// Response wraps the API response and provides methods to extract data.
type Response struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// raw is the parsed JSON response body.
	raw map[string]any
	// cfg is the compiled config for this resource.
	cfg *CompiledConfig
}

// Into unmarshals the response into a Terraform model struct.
func (r *Response) Into(ctx context.Context, model any) error {
	if r.raw == nil {
		return fmt.Errorf("no response data")
	}

	return Unmarshal(ctx, r.raw, model, r.cfg, nil)
}

// IntoWithPlan unmarshals the response into a model struct with plan preservation.
func (r *Response) IntoWithPlan(ctx context.Context, model any, plan any) error {
	if r.raw == nil {
		return fmt.Errorf("no response data")
	}

	return Unmarshal(ctx, r.raw, model, r.cfg, plan)
}

// ID extracts the resource ID from the response.
func (r *Response) ID() (int64, error) {
	if r.raw == nil {
		return 0, fmt.Errorf("no response data")
	}

	// Try to find ID in the response, handling envelope
	data := r.raw
	if r.cfg.raw.Envelope != nil && r.cfg.raw.Envelope.Response != "" {
		if wrapped, ok := data[r.cfg.raw.Envelope.Response]; ok {
			if m, ok := wrapped.(map[string]any); ok {
				data = m
			}
		}
	}

	if id, ok := data["id"]; ok {
		return toInt64(id)
	}

	// Search one level deep for a nested object containing "id"
	for _, val := range r.raw {
		if m, ok := val.(map[string]any); ok {
			if id, ok := m["id"]; ok {
				return toInt64(id)
			}
		}
	}

	return 0, fmt.Errorf("id not found in response")
}

// Raw returns the raw response map (for debugging or special cases).
func (r *Response) Raw() map[string]any {
	return r.raw
}

// Client is the main entry point for the translation layer.
// It holds compiled configs, the HTTP client, and version info.
type Client struct {
	httpClient *http.Client
	baseURL    string
	configs    map[string]*CompiledConfig
	version    string // Morpheus appliance version
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client (with auth transport).
func WithHTTPClient(c *http.Client) ClientOption {
	return func(cl *Client) {
		cl.httpClient = c
	}
}

// WithVersion sets the Morpheus appliance version for version-conditional transforms.
func WithVersion(ver string) ClientOption {
	return func(cl *Client) {
		cl.version = ver
	}
}

// NewClient creates a new translation layer client.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		configs:    make(map[string]*CompiledConfig),
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// RegisterResource adds a resource configuration to the client.
func (c *Client) RegisterResource(name string, configYAML []byte, opts ...ResourceOption) error {
	cfg, err := ParseConfig(configYAML)
	if err != nil {
		return fmt.Errorf("register resource %q: %w", name, err)
	}

	cc := Compile(cfg)

	for _, opt := range opts {
		opt(cc)
	}

	c.configs[name] = cc

	return nil
}

// RegisterResourceConfig adds a pre-parsed resource configuration to the client.
func (c *Client) RegisterResourceConfig(name string, cfg *ResourceConfig, opts ...ResourceOption) {
	cc := Compile(cfg)

	for _, opt := range opts {
		opt(cc)
	}

	c.configs[name] = cc
}

// ResourceOption configures a compiled resource config.
type ResourceOption func(*CompiledConfig)

// WithPostRead attaches a post-read hook to a resource config.
func WithPostRead(hook PostReadHook) ResourceOption {
	return func(cc *CompiledConfig) {
		cc.postRead = hook
	}
}

// WithPostWrite attaches a post-write hook to a resource config.
func WithPostWrite(hook PostWriteHook) ResourceOption {
	return func(cc *CompiledConfig) {
		cc.postWrite = hook
	}
}

// Execute performs a CRUD operation using the translation layer.
func (c *Client) Execute(ctx context.Context, req Request) (*Response, error) {
	cc, ok := c.configs[req.Resource]
	if !ok {
		return nil, fmt.Errorf("resource %q not registered", req.Resource)
	}

	// Apply version-specific overrides
	if c.version != "" {
		cc = cc.ResolveForVersion(c.version)
	}

	switch req.Operation {
	case Create:
		return c.doCreate(ctx, req, cc)
	case Read:
		return c.doRead(ctx, req, cc)
	case Update:
		return c.doUpdate(ctx, req, cc)
	case Delete:
		return c.doDelete(ctx, req, cc)
	case List:
		return c.doList(ctx, req, cc)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
	}
}

// GetConfig returns the compiled config for a resource (useful for testing).
func (c *Client) GetConfig(resource string) *CompiledConfig {
	return c.configs[resource]
}

// toInt64 converts a numeric value to int64.
func toInt64(v any) (int64, error) {
	switch n := v.(type) {
	case float64:
		return int64(n), nil
	case int:
		return int64(n), nil
	case int64:
		return n, nil
	case int32:
		return int64(n), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}
