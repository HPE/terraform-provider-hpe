// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	sdkv2client "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
)

// NewLegacyClient builds a client of the legacy (hand-written) SDK -- the one
// the sdkv2 resources use, constructed the way their provider constructs it --
// from the same TF_VAR_testacc_morpheus_* variables ProviderBlock reads, for a
// test that has to change something on the appliance behind Terraform's back
// with the sdkv2 resources' own request shapes. The generated SDK's client
// (NewClientForServer) is the first choice where its pruned API set has the
// operation.
func NewLegacyClient(t *testing.T) *morpheus.Client {
	t.Helper()

	url, ok := LookupProviderEnv("", "url")
	if !ok {
		t.Fatalf("%s not set", ProviderEnvName("", "url"))
	}

	insecure := false
	if v, ok := LookupProviderEnv("", "insecure"); ok {
		insecure = strings.EqualFold(strings.TrimSpace(v), "true")
	}

	token, _ := LookupProviderEnv("", "access_token")
	username, _ := LookupProviderEnv("", "username")
	password, _ := LookupProviderEnv("", "password")
	tenantSubdomain, _ := LookupProviderEnv("", "tenant_subdomain")

	if token == "" && (username == "" || password == "") {
		t.Fatal(errors.New("one of " + ProviderEnvName("", "access_token") + " or " +
			ProviderEnvName("", "username") + " and " + ProviderEnvName("", "password") + " must be set"))
	}

	return sdkv2client.NewLegacyClient(
		context.Background(), url, username, password, tenantSubdomain, token,
		morpheus.SkipLogin(), morpheus.WithInsecure(insecure),
	)
}

// LegacyError turns a legacy-SDK call's result into one error: the transport
// error when there is one, otherwise the appliance's refusal (status and the
// message in the body) when the response is not a success.
func LegacyError(resp *morpheus.Response, err error) error {
	if err != nil {
		return err
	}

	if resp == nil {
		return errors.New("no response")
	}

	if !resp.Success {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(resp.Body)))
	}

	return nil
}
