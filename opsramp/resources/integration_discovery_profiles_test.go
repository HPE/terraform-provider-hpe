// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIntegrationDiscoveryProfilesToAPI_OmitsEmptyResourceType(t *testing.T) {
	profiles := integrationDiscoveryProfilesToAPI([]IntegrationDiscoveryProfileModel{
		{
			Policy: &IntegrationDiscoveryPolicyModel{
				Rules: []IntegrationDiscoveryRuleModel{
					{
						FilterType: types.StringValue("ANY_CLOUD_RESOURCE"),
					},
				},
			},
		},
	})

	if len(profiles) != 1 || profiles[0].Policy == nil || len(profiles[0].Policy.Rules) != 1 {
		t.Fatalf("unexpected profile conversion result: %+v", profiles)
	}

	raw, err := json.Marshal(profiles[0].Policy.Rules[0])
	if err != nil {
		t.Fatalf("marshal rule: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal rule json: %v", err)
	}

	if _, exists := payload["resourceType"]; exists {
		t.Fatalf("resourceType should be omitted when unset or empty, got: %s", string(raw))
	}
}

func TestIntegrationDiscoveryProfilesToAPI_IncludesResourceTypeWhenPresent(t *testing.T) {
	profiles := integrationDiscoveryProfilesToAPI([]IntegrationDiscoveryProfileModel{
		{
			Policy: &IntegrationDiscoveryPolicyModel{
				Rules: []IntegrationDiscoveryRuleModel{
					{
						FilterType:   types.StringValue("ANY_CLOUD_RESOURCE"),
						ResourceType: []types.String{types.StringValue("VMWARE")},
					},
				},
			},
		},
	})

	if len(profiles) != 1 || profiles[0].Policy == nil || len(profiles[0].Policy.Rules) != 1 {
		t.Fatalf("unexpected profile conversion result: %+v", profiles)
	}

	raw, err := json.Marshal(profiles[0].Policy.Rules[0])
	if err != nil {
		t.Fatalf("marshal rule: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal rule json: %v", err)
	}

	resourceType, exists := payload["resourceType"]
	if !exists {
		t.Fatalf("resourceType should be present when values are set, got: %s", string(raw))
	}

	vals, ok := resourceType.([]any)
	if !ok || len(vals) != 1 || vals[0] != "VMWARE" {
		t.Fatalf("unexpected resourceType payload: %#v", resourceType)
	}
}
