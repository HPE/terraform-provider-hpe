// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingwhitelabel

import (
	"context"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func wlPtr[T any](v T) *T { return &v }

// TestMapResponseSpecAdditions verifies the read mapping of the MORPH-18721
// additions: security_banner, external_legal_text, supports_themes and the
// read-only account object.
func TestMapResponseSpecAdditions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	acct := sdk.ListWhitelabelSettings200ResponseWhitelabelSettingsAccount{
		Id:   wlPtr(int64(7)),
		Name: wlPtr("master"),
	}
	settings := &sdk.ListWhitelabelSettings200ResponseWhitelabelSettings{
		SecurityBanner:    wlPtr("banner"),
		ExternalLegalText: wlPtr("legal"),
		SupportsThemes:    wlPtr(true),
	}
	settings.Account.Set(&acct)

	var m SettingWhitelabelModel
	mapResponseToModel(ctx, &m, settings)

	if m.SecurityBanner.ValueString() != "banner" {
		t.Errorf("security_banner = %q, want banner", m.SecurityBanner.ValueString())
	}
	if m.ExternalLegalText.ValueString() != "legal" {
		t.Errorf("external_legal_text = %q, want legal", m.ExternalLegalText.ValueString())
	}
	if !m.SupportsThemes.ValueBool() {
		t.Errorf("supports_themes = %v, want true", m.SupportsThemes.ValueBool())
	}
	if m.Account.IsNull() || m.Account.Id.ValueInt64() != 7 || m.Account.Name.ValueString() != "master" {
		t.Errorf("account = %+v, want {id:7 name:master}", m.Account)
	}

	// Absent account maps to null.
	empty := &sdk.ListWhitelabelSettings200ResponseWhitelabelSettings{}
	var m2 SettingWhitelabelModel
	mapResponseToModel(ctx, &m2, empty)
	if !m2.Account.IsNull() {
		t.Errorf("expected account null for empty response, got %+v", m2.Account)
	}
}
