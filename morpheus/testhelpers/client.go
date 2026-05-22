// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

func newClient(ctx context.Context, t *testing.T) *sdk.APIClient {
	t.Helper()

	client, err := NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	return client
}

// NewClientForServer constructs an API client using the same TF_VAR naming
// convention as ProviderBlock. When preferredSystem is empty, it uses
// the non-server-specific test variable names.
func NewClientForServer(ctx context.Context, preferredSystem string) (*sdk.APIClient, error) {
	var username, password string

	url, ok := LookupProviderEnv(preferredSystem, "url")
	if !ok {
		return nil, fmt.Errorf("%s not set", ProviderEnvName(preferredSystem, "url"))
	}

	token, ok := LookupProviderEnv(preferredSystem, "access_token")
	if !ok {
		username, ok = LookupProviderEnv(preferredSystem, "username")
		if !ok {
			return nil, errors.New(
				"one of " + ProviderEnvName(preferredSystem, "access_token") + " or " +
					ProviderEnvName(preferredSystem, "username") + " must be set",
			)
		}

		password, ok = LookupProviderEnv(preferredSystem, "password")
		if !ok {
			return nil, errors.New(
				"one of " + ProviderEnvName(preferredSystem, "access_token") + " or " +
					ProviderEnvName(preferredSystem, "password") + " must be set",
			)
		}
	}

	tenantSubdomain, _ := LookupProviderEnv(preferredSystem, "tenant_subdomain")

	_, insecure := LookupProviderEnv(preferredSystem, "insecure")
	var opts []clientfactory.ClientOption
	if insecure {
		opts = append(opts, clientfactory.WithInsecureTLS())
	}

	client := clientfactory.NewAPIClient(
		ctx,
		url,
		username,
		password,
		tenantSubdomain,
		token,
		opts...,
	)

	return client, nil
}

// LookupProviderEnv looks up a TF_VAR provider input using the same naming
// convention as ProviderBlock, falling back to the
// shared non-server-specific name when a server-specific value is absent.
func LookupProviderEnv(preferredSystem string, suffix string) (string, bool) {
	if name := ProviderEnvName(preferredSystem, suffix); name != "" {
		if value, ok := os.LookupEnv(name); ok {
			return value, true
		}
	}

	return os.LookupEnv("TF_VAR_testacc_morpheus_" + suffix)
}

// ProviderEnvName returns the TF_VAR name used by the acceptance-test provider
// config for the given server-specific suffix.
func ProviderEnvName(preferredSystem string, suffix string) string {
	preferredSystem = strings.TrimSpace(preferredSystem)
	if preferredSystem == "" || strings.EqualFold(preferredSystem, "all") {
		return "TF_VAR_testacc_morpheus_" + suffix
	}

	return fmt.Sprintf("TF_VAR_testacc_morpheus_%s_%s", preferredSystem, suffix)
}
