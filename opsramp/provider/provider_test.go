// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provider_test

import (
	"context"
	"testing"

	providerimpl "github.com/HPE/terraform-provider-hpe/opsramp/provider"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
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
