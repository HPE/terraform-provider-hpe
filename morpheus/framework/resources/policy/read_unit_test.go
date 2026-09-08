// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package policy

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMotdTitleToState guards the read-path nil-safety fix for the MOTD policy title: an
// explicit JSON null (IsSet true, nil Get) or an omitted title must map to a null string
// without dereferencing a nil pointer. A present title maps through.
func TestMotdTitleToState(t *testing.T) {
	t.Parallel()

	if got := motdTitleToState(*sdk.NewNullableString(nil)); !got.IsNull() {
		t.Errorf("explicit null title = %v, want null", got)
	}

	if got := motdTitleToState(sdk.NullableString{}); !got.IsNull() {
		t.Errorf("absent title = %v, want null", got)
	}

	title := "Scheduled maintenance"
	if got := motdTitleToState(*sdk.NewNullableString(&title)); got.IsNull() || got.ValueString() != title {
		t.Errorf("present title = %v, want %q", got, title)
	}
}
