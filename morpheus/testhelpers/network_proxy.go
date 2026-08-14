// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"crypto/rand"
	"net/http"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// CreateNetworkProxy creates a network proxy fixture via the API and returns its
// id and name.
//
// There is no hpe_morpheus_network_proxy resource, so the network proxy data
// source test builds its own fixture through the SDK rather than depending on a
// proxy that happens to exist on the appliance. This keeps the test
// self-contained and free of environment configuration.
func CreateNetworkProxy(t *testing.T) (int64, string) {
	t.Helper()

	ctx := context.TODO()
	client := newClient(ctx, t)

	name := "tfacc-proxy-" + rand.Text()
	host := "127.0.0.1"
	port := "3128"
	visibility := "private"

	req := sdk.CreateNetworkProxyRequest{
		NetworkProxy: &sdk.CreateNetworkProxyRequestNetworkProxy{
			Name:       &name,
			ProxyHost:  &host,
			ProxyPort:  &port,
			Visibility: &visibility,
		},
	}

	resp, hresp, err := client.NetworksAPI.CreateNetworkProxy(ctx).
		CreateNetworkProxyRequest(req).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK ||
		resp == nil || resp.NetworkProxy == nil || resp.NetworkProxy.Id == nil {
		t.Fatalf("failed to create network proxy fixture %q: %v", name, err)
	}

	id := *resp.NetworkProxy.Id
	t.Logf("created network proxy fixture %d (%s)", id, name)

	return id, name
}

// DeleteNetworkProxy deletes a network proxy fixture created by
// CreateNetworkProxy. A missing proxy (already deleted) is not an error.
func DeleteNetworkProxy(t *testing.T, id int64) {
	t.Helper()

	ctx := context.TODO()
	client := newClient(ctx, t)

	_, hresp, err := client.NetworksAPI.DeleteNetworkProxy(ctx, id).Execute()
	if err != nil && (hresp == nil || hresp.StatusCode != http.StatusNotFound) {
		t.Logf("warning: failed to delete network proxy fixture %d: %v", id, err)

		return
	}

	t.Logf("deleted network proxy fixture %d", id)
}
