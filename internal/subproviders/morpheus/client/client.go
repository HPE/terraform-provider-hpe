// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package client

import (
	"net/http"
)

type Client struct {
	HTTPClient *http.Client // HTTP client to use for requests
	URL        string
}
