// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A morpheus block that carries no url became valid config once url was made
// optional, so that the framework provider can obtain the connection details
// from a greenlake_connected block instead. The legacy provider is configured
// with the same config and must not panic on it.
func TestProviderConfigureWithoutURL(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		morpheus []interface{}
		wantMeta string
	}{
		{
			name:     "empty block",
			morpheus: []interface{}{nil},
			wantMeta: incompleteMorpheusBlock,
		},
		{
			name: "greenlake_connected only",
			morpheus: []interface{}{
				map[string]interface{}{
					"greenlake_connected": []interface{}{
						map[string]interface{}{"client_id": "id"},
					},
				},
			},
			wantMeta: incompleteMorpheusBlock,
		},
		{
			name:     "no block",
			morpheus: nil,
			wantMeta: missingMorpheusBlock,
		},
		{
			name:     "empty list",
			morpheus: []interface{}{},
			wantMeta: missingMorpheusBlock,
		},
	}

	for _, testcase := range testcases {
		tc := testcase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(
				t,
				Provider().Schema,
				map[string]interface{}{},
			)

			if tc.morpheus != nil {
				if err := d.Set("morpheus", tc.morpheus); err != nil {
					t.Fatalf("failed to set morpheus block: %v", err)
				}
			}

			// Must not panic.
			meta, diags := providerConfigure(context.Background(), d)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			msg, ok := meta.(string)
			if !ok {
				t.Fatalf("expected a guidance message, got %T", meta)
			}

			if !strings.EqualFold(msg, tc.wantMeta) {
				t.Errorf("unexpected guidance message:\ngot:\n%s\nwant:\n%s", msg, tc.wantMeta)
			}
		})
	}
}

// A block with a url must still produce a working client.
func TestProviderConfigureWithURL(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

	err := d.Set("morpheus", []interface{}{
		map[string]interface{}{
			"url":          "https://example.com",
			"access_token": "token",
		},
	})
	if err != nil {
		t.Fatalf("failed to set morpheus block: %v", err)
	}

	meta, diags := providerConfigure(context.Background(), d)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if _, ok := meta.(string); ok {
		t.Fatalf("expected a client, got guidance message: %v", meta)
	}
}
