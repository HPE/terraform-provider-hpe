// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestMuxedProviderSchemaHasNoConflicts builds the real tf6muxserver the same
// way main.go and the acceptance harness do (the plugin-framework server and the
// SDKv2 server behind the mux), then calls GetProviderSchema.
//
// The mux reports a type name registered in both servers as an error-severity
// diagnostic on that RPC ("Invalid Provider Server Combination ... Duplicate
// resource type: ..."), returning a nil Go error. That is the failure mode a
// SDKv2 -> framework port such as hpe_morpheus_tenant must avoid, and it
// surfaces at terraform plan rather than at build or startup.
//
// This runs in plain `go test` with no TF_ACC or appliance (GetProviderSchema
// never configures the provider or performs network calls) and covers every
// type namespace the mux checks (resources, data sources, ephemeral resources,
// functions, list resources, state stores, actions).
func TestMuxedProviderSchemaHasNoConflicts(t *testing.T) {
	factory := testhelpers.GetAccTestFactories(
		t, adapter.NewMorpheus(), sdkv2morpheus.Provider(),
	)["hpe"]

	server, err := factory()
	if err != nil {
		t.Fatalf("building muxed provider server: %s", err)
	}

	resp, err := server.GetProviderSchema(
		context.Background(), &tfprotov6.GetProviderSchemaRequest{},
	)
	if err != nil {
		t.Fatalf("GetProviderSchema: %s", err)
	}

	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Errorf("muxed provider schema error: %s: %s", d.Summary, d.Detail)
		}
	}
}
