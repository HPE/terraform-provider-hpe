// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provider_test

import (
	"context"
	"testing"

	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	providerimpl "github.com/HPE/terraform-provider-hpe/opsramp/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := providerimpl.New("test")()

	var resp frameworkprovider.MetadataResponse
	p.Metadata(context.Background(), frameworkprovider.MetadataRequest{}, &resp)

	if resp.TypeName != "opsramp" {
		t.Fatalf("expected provider type name %q, got %q", "opsramp", resp.TypeName)
	}

	if resp.Version != "test" {
		t.Fatalf("expected provider version %q, got %q", "test", resp.Version)
	}
}

func TestProviderSchemaHasCoreAttributes(t *testing.T) {
	p := providerimpl.New("test")()

	var resp frameworkprovider.SchemaResponse
	p.Schema(context.Background(), frameworkprovider.SchemaRequest{}, &resp)

	requiredAttrs := []string{"client_id", "client_secret", "endpoint", "tenant"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("expected provider schema to include %q attribute", attr)
		}
	}
}

func TestProviderProtocolSchemaIsValid(t *testing.T) {
	server, err := providerserver.NewProtocol6WithError(providerimpl.New("test")())()
	if err != nil {
		t.Fatalf("could not create provider server: %v", err)
	}

	response, err := server.GetProviderSchema(
		context.Background(),
		&tfprotov6.GetProviderSchemaRequest{},
	)
	if err != nil {
		t.Fatalf("GetProviderSchema() unexpected error: %v", err)
	}

	for _, diagnostic := range response.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("GetProviderSchema() error: %s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
}
