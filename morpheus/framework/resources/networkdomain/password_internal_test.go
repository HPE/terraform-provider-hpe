// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// domainPasswordForUpdate must ALWAYS return a value (never nil): omitting
// domainPassword on update clears the stored password. It sends the configured
// value only when re-sending, the deprecated plaintext attribute whenever it is
// set, and the mask otherwise.
func TestUnitDomainPasswordForUpdate(t *testing.T) {
	t.Parallel()

	secret := "Secret123!"

	legacy := "Legacy123!"
	none := types.StringNull()

	cases := []struct {
		name   string
		resend bool
		wo     types.String
		plain  types.String
		want   string
	}{
		{"resend with configured value sends the secret", true, types.StringValue(secret), none, secret},
		{"resend but no configured value sends the mask", true, none, none, passwordMask},
		{"no resend sends the mask even when configured", false, types.StringValue(secret), none, passwordMask},
		{"no resend and no value sends the mask", false, none, none, passwordMask},
		{"unknown value is treated as unset", true, types.StringUnknown(), none, passwordMask},
		// The deprecated plaintext attribute is sent on every update while set,
		// as the SDKv2 resource did.
		{"deprecated domain_password is sent without resend", false, none, types.StringValue(legacy), legacy},
		{"deprecated domain_password is sent with resend", true, none, types.StringValue(legacy), legacy},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := domainPasswordForUpdate(tc.resend, tc.wo, tc.plain)
			if got == nil {
				t.Fatal("domainPasswordForUpdate must never return nil: omitting the field clears the password")
			}

			if *got != tc.want {
				t.Errorf("got %q, want %q", *got, tc.want)
			}
		})
	}
}
