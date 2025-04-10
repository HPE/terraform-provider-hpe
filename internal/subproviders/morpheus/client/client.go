// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package client

import (
	"net/http"
)

type Client struct {
	// For now just use a standard http client
	http.Client
}
