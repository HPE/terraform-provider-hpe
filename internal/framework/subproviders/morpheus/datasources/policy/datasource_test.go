// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestNewDataSource(t *testing.T) {
	t.Parallel()

	ds := NewDataSource()
	if ds == nil {
		t.Fatal("NewDataSource returned nil")
	}

	// Test that it implements the DataSource interface
	var _ datasource.DataSource = ds
}

func TestDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := NewDataSource()
	req := datasource.MetadataRequest{
		ProviderTypeName: "hpe",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "hpe_morpheus_policy" {
		t.Errorf("Expected TypeName to be 'hpe_morpheus_policy', got '%s'", resp.TypeName)
	}
}

func TestDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := NewDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema attributes should not be nil")
	}

	// Check that required attributes exist
	if _, exists := resp.Schema.Attributes["id"]; !exists {
		t.Error("Schema should have 'id' attribute")
	}

	if _, exists := resp.Schema.Attributes["name"]; !exists {
		t.Error("Schema should have 'name' attribute")
	}

	if _, exists := resp.Schema.Attributes["policy_type"]; !exists {
		t.Error("Schema should have 'policy_type' attribute")
	}

	if _, exists := resp.Schema.Attributes["enabled"]; !exists {
		t.Error("Schema should have 'enabled' attribute")
	}
}
