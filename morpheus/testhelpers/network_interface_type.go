// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"net/http"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// NetworkInterfaceTypeFixture identifies a network interface (NIC) type that
// exists on the appliance under test, discovered at runtime so the test needs
// no environment configuration or hard-coded ids.
type NetworkInterfaceTypeFixture struct {
	CloudID           int64
	ProvisionTypeCode string
	Name              string
	ID                int64
}

// nicProvisionTypeCandidates are the provision types whose zoneNetworkOptions
// carry NIC types, in preference order. vmware is first to match the data
// source's default and exercise the primary path where possible.
var nicProvisionTypeCandidates = []string{"vmware", "kvm", "esxi", "hyperv", "xen"}

// DiscoverNetworkInterfaceType finds a (cloud, provision type, NIC type) triple
// on the appliance for which the network interface type lookup resolves a
// concrete id. It skips the test when no cloud yields NIC types.
//
// A NIC type is scoped to a cloud, which cannot be created from Terraform
// without provider credentials, so the fixture is discovered from the appliance
// (mirroring CreateInstance) rather than read from the environment.
func DiscoverNetworkInterfaceType(t *testing.T) NetworkInterfaceTypeFixture {
	t.Helper()

	ctx := context.TODO()
	client := newClient(ctx, t)

	zonesResp, hresp, err := client.CloudsAPI.ListClouds(ctx).Max(200).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK || zonesResp == nil {
		t.Fatalf("failed to list clouds for NIC type discovery: %v", err)
	}

	for _, code := range nicProvisionTypeCandidates {
		ptID, ok := provisionTypeIDByCode(ctx, client, code)
		if !ok {
			continue
		}

		for i := range zonesResp.Zones {
			z := &zonesResp.Zones[i]
			if z.Id == nil {
				continue
			}

			opts, oResp, oErr := client.OptionsAPI.ListOptionNetworkOptions(ctx).
				ZoneId(*z.Id).ProvisionTypeId(ptID).Execute()
			if oErr != nil || oResp == nil || oResp.StatusCode != http.StatusOK ||
				opts == nil || opts.Data == nil {
				continue
			}

			for j := range opts.Data.NetworkTypes {
				nt := &opts.Data.NetworkTypes[j]
				if nt.Id != nil && nt.Name != nil {
					fixture := NetworkInterfaceTypeFixture{
						CloudID:           *z.Id,
						ProvisionTypeCode: code,
						Name:              *nt.Name,
						ID:                *nt.Id,
					}
					t.Logf("discovered NIC type fixture: cloud=%d provision_type=%s name=%q id=%d",
						fixture.CloudID, fixture.ProvisionTypeCode, fixture.Name, fixture.ID)

					return fixture
				}
			}
		}
	}

	t.Skip("no cloud on the appliance returns network interface types; skipping")

	return NetworkInterfaceTypeFixture{}
}

// provisionTypeIDByCode resolves a provision type code to its id, returning
// false when the code is not present on the appliance.
func provisionTypeIDByCode(
	ctx context.Context,
	client *sdk.APIClient,
	code string,
) (int64, bool) {
	resp, hresp, err := client.ProvisioningAPI.ListProvisionTypes(ctx).Code(code).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK || resp == nil {
		return 0, false
	}

	for i := range resp.ProvisionTypes {
		pt := &resp.ProvisionTypes[i]
		if pt.Code != nil && *pt.Code == code && pt.Id != nil {
			return *pt.Id, true
		}
	}

	return 0, false
}
