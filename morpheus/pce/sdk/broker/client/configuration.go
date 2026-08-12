// (C) Copyright 2021-2026 Hewlett Packard Enterprise Development LP

package client

import (
	"net/http"
)

// contextKey identifies the type of value in the context. Since these are
// strings, a short description is available for logging and debugging.
type contextKey string

func (c contextKey) String() string {
	return "auth " + string(c)
}

var (
	// ContextBasicAuth takes BasicAuth as authentication for the request.
	ContextBasicAuth = contextKey("basic")

	// ContextAccessToken takes a string oauth2 access token as authentication for the request.
	ContextAccessToken = contextKey("accesstoken")
)

// BasicAuth provides basic http authentication to a request passed via context using ContextBasicAuth
type BasicAuth struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

type Configuration struct {
	Host               string            `json:"host,omitempty"`
	Scheme             string            `json:"scheme,omitempty"`
	DefaultHeader      map[string]string `json:"defaultHeader,omitempty"`
	DefaultQueryParams map[string]string `json:"defaultQueryParams"`
	UserAgent          string            `json:"userAgent,omitempty"`
	HTTPClient         *http.Client
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{
		DefaultHeader:      make(map[string]string),
		DefaultQueryParams: make(map[string]string),
		UserAgent:          "vmmas/cmp/go-sdk",
	}

	return cfg
}

func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}

func (c *Configuration) AddDefaultQueryParams(key string, value string) {
	c.DefaultQueryParams[key] = value
}
