// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func TestIsMetal_CorrectCode(t *testing.T) {
	t.Parallel()

	code := metalProvisionTypeCode
	inst := &sdk.GetInstance200ResponseInstance{
		Layout: &sdk.GetInstance200ResponseInstanceLayout{
			ProvisionTypeCode: &code,
		},
	}

	if !IsMetal(inst) {
		t.Error("expected IsMetal=true for HPE bare-metal provision type")
	}
}

func TestIsMetal_NonMetal(t *testing.T) {
	t.Parallel()

	code := "vmware"
	inst := &sdk.GetInstance200ResponseInstance{
		Layout: &sdk.GetInstance200ResponseInstanceLayout{
			ProvisionTypeCode: &code,
		},
	}

	if IsMetal(inst) {
		t.Error("expected IsMetal=false for vmware provision type")
	}
}

func TestIsMetal_NilLayout(t *testing.T) {
	t.Parallel()

	inst := &sdk.GetInstance200ResponseInstance{
		Layout: nil,
	}

	if IsMetal(inst) {
		t.Error("expected IsMetal=false when layout is nil")
	}
}

func TestIsMetal_NilProvisionTypeCode(t *testing.T) {
	t.Parallel()

	inst := &sdk.GetInstance200ResponseInstance{
		Layout: &sdk.GetInstance200ResponseInstanceLayout{
			ProvisionTypeCode: nil,
		},
	}

	if IsMetal(inst) {
		t.Error("expected IsMetal=false when provisionTypeCode is nil")
	}
}

func TestIsMetal_NilInstance(t *testing.T) {
	t.Parallel()

	if IsMetal(nil) {
		t.Error("expected IsMetal=false for nil instance")
	}
}

// format=="bareMetal" with a non-HPE provision type must be REJECTED.
// This tests that keying on format would be wrong — we key on
// provisionTypeCode, not format.
func TestIsMetal_BareMetalFormatNonHPE(t *testing.T) {
	t.Parallel()

	// An instance with format=bareMetal but a non-HPE provision type
	// (e.g. OneView, MaaS) must NOT pass the metal check.
	code := "oneview-provision"
	inst := &sdk.GetInstance200ResponseInstance{
		Layout: &sdk.GetInstance200ResponseInstanceLayout{
			ProvisionTypeCode: &code,
		},
	}

	if IsMetal(inst) {
		t.Error("expected IsMetal=false for oneview-provision; " +
			"format==bareMetal does not imply HPE metal")
	}
}

func TestProvisionTypeCode_Absent(t *testing.T) {
	t.Parallel()

	inst := &sdk.GetInstance200ResponseInstance{
		Layout: &sdk.GetInstance200ResponseInstanceLayout{},
	}

	code, ok := provisionTypeCode(inst)
	if ok {
		t.Errorf("expected ok=false, got code=%q", code)
	}
}

func TestNotMetalDetail_ContainsBothCodes(t *testing.T) {
	t.Parallel()

	detail := notMetalDetail(42, "vmware")
	if detail == "" {
		t.Fatal("expected non-empty detail")
	}

	// Must mention both the metal code and the actual code.
	for _, want := range []string{metalProvisionTypeCode, "vmware", "42"} {
		if !contains(detail, want) {
			t.Errorf("detail %q does not contain %q", detail, want)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
