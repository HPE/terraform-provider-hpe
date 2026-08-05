// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildAddNodeEnvelope_PreProvisionedFalse(t *testing.T) {
	t.Parallel()

	plan := &instanceNodeModel{
		ResourcePoolID: types.Int64Value(42),
		PreProvisioned: types.BoolValue(false),
	}

	env := buildAddNodeEnvelope(plan, "generic-add-node")

	// preProvisioned false/null => key omitted entirely, never "off"
	if _, ok := env["preProvisioned"]; ok {
		t.Error("preProvisioned should be omitted when false, not sent as 'off'")
	}
}

func TestBuildAddNodeEnvelope_PreProvisionedNull(t *testing.T) {
	t.Parallel()

	plan := &instanceNodeModel{
		ResourcePoolID: types.Int64Value(42),
		PreProvisioned: types.BoolNull(),
	}

	env := buildAddNodeEnvelope(plan, "generic-add-node")

	if _, ok := env["preProvisioned"]; ok {
		t.Error("preProvisioned should be omitted when null")
	}
}

func TestBuildAddNodeEnvelope_PreProvisionedTrue(t *testing.T) {
	t.Parallel()

	plan := &instanceNodeModel{
		ResourcePoolID:   types.Int64Value(42),
		PreProvisioned:   types.BoolValue(true),
		SelectedServerID: types.Int64Value(99),
		SshHost:          types.StringValue("10.0.0.1"),
		SshUsername:      types.StringValue("root"),
		SshPassword:      types.StringValue("secret"),
		SshKeyPairID:     types.Int64Value(7),
	}

	env := buildAddNodeEnvelope(plan, "generic-add-node")

	if env["preProvisioned"] != "on" {
		t.Errorf("expected preProvisioned='on', got %v", env["preProvisioned"])
	}

	if env["selectedServerId"] != int64(99) {
		t.Errorf("expected selectedServerId=99, got %v", env["selectedServerId"])
	}

	if env["sshHost"] != "10.0.0.1" {
		t.Errorf("expected sshHost='10.0.0.1', got %v", env["sshHost"])
	}
}

func TestBuildAddNodeEnvelope_AlwaysSendsCount1(t *testing.T) {
	t.Parallel()

	plan := &instanceNodeModel{
		ResourcePoolID: types.Int64Value(1),
		PreProvisioned: types.BoolNull(),
	}

	env := buildAddNodeEnvelope(plan, "generic-add-node")

	if env["count"] != 1 {
		t.Errorf("expected count=1, got %v", env["count"])
	}
}

func TestBuildAddNodeEnvelope_PoolFormat(t *testing.T) {
	t.Parallel()

	plan := &instanceNodeModel{
		ResourcePoolID: types.Int64Value(42),
		PreProvisioned: types.BoolNull(),
	}

	env := buildAddNodeEnvelope(plan, "generic-add-node")

	if env["selectedResourcePoolId"] != "pool-42" {
		t.Errorf("expected selectedResourcePoolId='pool-42', got %v",
			env["selectedResourcePoolId"])
	}
}

func TestBuildAddNodeEnvelope_NoPoolOmitsKey(t *testing.T) {
	t.Parallel()

	plan := &instanceNodeModel{
		ResourcePoolID: types.Int64Null(),
		PreProvisioned: types.BoolNull(),
	}

	env := buildAddNodeEnvelope(plan, "generic-add-node")

	if _, ok := env["selectedResourcePoolId"]; ok {
		t.Error("selectedResourcePoolId should be absent when resource_pool_id is null")
	}

	// count and code must still be present.
	if env["count"] != 1 {
		t.Errorf("expected count=1, got %v", env["count"])
	}

	if env["code"] != "generic-add-node" {
		t.Errorf("expected code='generic-add-node', got %v", env["code"])
	}
}

func TestExtractContainerIDFromResults_Valid(t *testing.T) {
	t.Parallel()

	results := map[string]any{
		"1417": map[string]any{
			"containers": []any{
				map[string]any{"id": float64(10792)},
			},
		},
	}

	id, err := extractContainerIDFromResults(results, 1417)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if id != 10792 {
		t.Errorf("expected 10792, got %d", id)
	}
}

func TestExtractContainerIDFromResults_NoContainers(t *testing.T) {
	t.Parallel()

	results := map[string]any{
		"1417": map[string]any{
			"containers": []any{},
		},
	}

	_, err := extractContainerIDFromResults(results, 1417)
	if err == nil {
		t.Error("expected error for empty containers")
	}
}

func TestExtractContainerIDFromResults_WrongInstance(t *testing.T) {
	t.Parallel()

	results := map[string]any{
		"999": map[string]any{},
	}

	_, err := extractContainerIDFromResults(results, 1417)
	if err == nil {
		t.Error("expected error for missing instance key")
	}
}
