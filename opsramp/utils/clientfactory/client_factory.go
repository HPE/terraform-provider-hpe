// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"fmt"
	"sync"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

// ClientFactory stores the OpsRamp provider configuration and creates the
// API client lazily on first use. This avoids blocking provider Configure
// with network calls (OAuth token retrieval), which would prevent subsequent
// Terraform operations from working if the initial call fails.
type ClientFactory struct {
	clientID     string
	clientSecret string
	endpoint     string
	tenant       string

	mu     sync.Mutex
	client *client.OpsRampClient
}

// NewClientFactory creates a new ClientFactory with the given configuration.
// No network calls are made at this point.
func NewClientFactory(clientID, clientSecret, endpoint, tenant string) *ClientFactory {
	return &ClientFactory{
		clientID:     clientID,
		clientSecret: clientSecret,
		endpoint:     endpoint,
		tenant:       tenant,
	}
}

// Client returns the OpsRamp API client, creating it on first successful call.
// Successful results are cached; failures are retried on subsequent calls.
func (f *ClientFactory) Client() (*client.OpsRampClient, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.client != nil {
		return f.client, nil
	}

	c, err := client.NewOpsRampClient(f.clientID, f.clientSecret, f.endpoint, f.tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpsRamp client: %w", err)
	}

	f.client = c

	return f.client, nil
}
