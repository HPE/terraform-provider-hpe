// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strp(s string) *string { return &s }
func i64p(i int64) *int64   { return &i }

func TestFindInstalledLicense(t *testing.T) {
	const hash = "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"
	upperHash := "AABBCCDDEEFF00112233445566778899AABBCCDDEEFF00112233445566778899"

	installed := sdk.GetLicense200ResponseInstalledLicensesInner{
		KeyId: strp("INST1234"),
		Hash:  strp(hash),
	}
	otherInstalled := sdk.GetLicense200ResponseInstalledLicensesInner{
		KeyId: strp("OTHER567"),
		Hash:  strp("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	}

	t.Run("matches by hash case-insensitively", func(t *testing.T) {
		result := &sdk.GetLicense200Response{
			InstalledLicenses: []sdk.GetLicense200ResponseInstalledLicensesInner{otherInstalled, installed},
		}
		got := findInstalledLicense(result, upperHash, "")
		if got == nil || got.KeyId != "INST1234" {
			t.Fatalf("expected INST1234 by hash, got %+v", got)
		}
	})

	t.Run("hash takes precedence over keyId", func(t *testing.T) {
		// keyID points at otherInstalled, but hash points at installed: hash wins.
		result := &sdk.GetLicense200Response{
			InstalledLicenses: []sdk.GetLicense200ResponseInstalledLicensesInner{otherInstalled, installed},
		}
		got := findInstalledLicense(result, hash, "OTHER567")
		if got == nil || got.KeyId != "INST1234" {
			t.Fatalf("expected hash match INST1234, got %+v", got)
		}
	})

	t.Run("matches by keyId when hash empty", func(t *testing.T) {
		result := &sdk.GetLicense200Response{
			InstalledLicenses: []sdk.GetLicense200ResponseInstalledLicensesInner{otherInstalled, installed},
		}
		got := findInstalledLicense(result, "", "INST1234")
		if got == nil || got.KeyId != "INST1234" {
			t.Fatalf("expected INST1234 by keyId, got %+v", got)
		}
	})

	t.Run("InstalledLicenses preferred over License", func(t *testing.T) {
		result := &sdk.GetLicense200Response{
			InstalledLicenses: []sdk.GetLicense200ResponseInstalledLicensesInner{installed},
			License: &sdk.GetLicense200ResponseLicense{
				KeyId: strp("MERGED12"),
				Hash:  strp(hash),
			},
		}
		got := findInstalledLicense(result, hash, "")
		if got == nil || got.KeyId != "INST1234" {
			t.Fatalf("expected installed INST1234 over merged, got %+v", got)
		}
	})

	t.Run("falls back to License when list has no match", func(t *testing.T) {
		result := &sdk.GetLicense200Response{
			InstalledLicenses: []sdk.GetLicense200ResponseInstalledLicensesInner{otherInstalled},
			License: &sdk.GetLicense200ResponseLicense{
				KeyId: strp("MERGED12"),
				Hash:  strp(hash),
			},
		}
		got := findInstalledLicense(result, hash, "")
		if got == nil || got.KeyId != "MERGED12" {
			t.Fatalf("expected merged MERGED12, got %+v", got)
		}
	})

	t.Run("nil result", func(t *testing.T) {
		if got := findInstalledLicense(nil, hash, ""); got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		if got := findInstalledLicense(&sdk.GetLicense200Response{}, hash, "INST1234"); got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})

	t.Run("nothing to match on", func(t *testing.T) {
		result := &sdk.GetLicense200Response{
			InstalledLicenses: []sdk.GetLicense200ResponseInstalledLicensesInner{installed},
		}
		if got := findInstalledLicense(result, "", ""); got != nil {
			t.Fatalf("expected nil when no hash/keyId, got %+v", got)
		}
	})
}

func TestLicenseToState(t *testing.T) {
	t.Run("nil end date yields null", func(t *testing.T) {
		start := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
		st := licenseToState(installedLicense{
			KeyId:     "KEY12345",
			StartDate: &start,
			EndDate:   nil,
		})
		if !st.EndDate.IsNull() {
			t.Fatalf("expected EndDate null, got %v", st.EndDate)
		}
		if st.StartDate.ValueString() != "2026-01-02T03:04:05Z" {
			t.Fatalf("StartDate = %q", st.StartDate.ValueString())
		}
		if st.Id.ValueString() != "KEY12345" {
			t.Fatalf("Id = %q", st.Id.ValueString())
		}
		if !st.Key.IsNull() {
			t.Fatalf("Key must be null (write-only)")
		}
	})

	t.Run("end date yields RFC3339 UTC", func(t *testing.T) {
		end := time.Date(2027, 12, 31, 23, 0, 0, 0, time.UTC)
		st := licenseToState(installedLicense{EndDate: &end})
		if st.EndDate.ValueString() != "2027-12-31T23:00:00Z" {
			t.Fatalf("EndDate = %q", st.EndDate.ValueString())
		}
	})

	t.Run("nil ints yield null", func(t *testing.T) {
		st := licenseToState(installedLicense{})
		if !st.MaxInstances.IsNull() || !st.MaxMemory.IsNull() || !st.MaxStorage.IsNull() {
			t.Fatalf("expected nil ints to be null")
		}
	})

	t.Run("set ints carried", func(t *testing.T) {
		st := licenseToState(installedLicense{MaxInstances: i64p(5)})
		if st.MaxInstances.ValueInt64() != 5 {
			t.Fatalf("MaxInstances = %d", st.MaxInstances.ValueInt64())
		}
	})
}

// TestProvisionalState: the partial state left behind when the key was
// installed but not read back holds the key id and nothing else, so that
// Read and Delete can identify the licence and the next refresh completes the
// details. Checked through the Terraform value so a detail added to the model
// cannot slip through non-null.
func TestProvisionalState(t *testing.T) {
	ctx := context.Background()

	st := provisionalState("KEY12345")

	var s tfsdk.State
	s.Schema = LicenseResourceSchema(ctx)
	s.Raw = tftypes.NewValue(s.Schema.Type().TerraformType(ctx), nil)

	if diags := s.Set(ctx, st); diags.HasError() {
		t.Fatalf("provisional state does not fit the schema: %v", diags)
	}

	var attrs map[string]tftypes.Value
	if err := s.Raw.As(&attrs); err != nil {
		t.Fatal(err)
	}

	for name, v := range attrs {
		switch name {
		case "id":
			var id string
			if err := v.As(&id); err != nil || id != "KEY12345" {
				t.Fatalf("id = %v, want KEY12345 (err %v)", v, err)
			}
		default:
			if !v.IsNull() {
				t.Errorf("%s = %v, want null", name, v)
			}
		}
	}
}
