// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitBuildVmwareCloudConfigPassword pins the write-only contract for
// config_vmware.password.
//
// Because the attribute is declared WriteOnly, Terraform strips its value from
// the plan (and state), leaving it null there - the value is only present in
// the request config. buildVmwareCloudConfig must therefore source the password
// from configModel (req.Config), never the plan. Sourcing it from the plan was
// the regression that dropped the password and made VMware cloud creation fail
// with HTTP 400 {"errors":{"password":"Enter your password"}}.
func TestUnitBuildVmwareCloudConfigPassword(t *testing.T) {
	t.Parallel()

	const secret = "s3cret-from-config"

	tests := []struct {
		name         string
		planPassword types.String
		cfgPassword  types.String
		wantPassword string
		wantSet      bool
	}{
		{
			// The real write-only scenario: null in the plan, set in config.
			// A plan-sourced read would yield nil here and fail.
			name:         "sourced from config when plan is null",
			planPassword: types.StringNull(),
			cfgPassword:  types.StringValue(secret),
			wantPassword: secret,
			wantSet:      true,
		},
		{
			// Config always wins; guards against re-introducing a plan read
			// even if a stale value somehow appears in the plan.
			name:         "config wins over any plan value",
			planPassword: types.StringValue("stale-plan-value"),
			cfgPassword:  types.StringValue(secret),
			wantPassword: secret,
			wantSet:      true,
		},
		{
			// No password in config -> nothing sent to the API.
			name:         "unset when config is null",
			planPassword: types.StringNull(),
			cfgPassword:  types.StringNull(),
			wantSet:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			plan := CloudModel{
				ConfigVmware: ConfigVmwareValue{
					ApiUrl:     types.StringValue("https://vcenter.example.local"),
					ApiVersion: types.StringValue("7.0"),
					Datacenter: types.StringValue("DC9"),
					Username:   types.StringValue("administrator@vsphere.local"),
					Password:   tc.planPassword,
					state:      attr.ValueStateKnown,
				},
			}
			configModel := CloudModel{
				ConfigVmware: ConfigVmwareValue{
					Password: tc.cfgPassword,
					state:    attr.ValueStateKnown,
				},
			}

			got := buildVmwareCloudConfig(plan, configModel)

			if tc.wantSet {
				if got.Password == nil {
					t.Fatalf("password: got nil, want %q (write-only value must come from req.Config)", tc.wantPassword)
				}
				if *got.Password != tc.wantPassword {
					t.Fatalf("password: got %q, want %q", *got.Password, tc.wantPassword)
				}
			} else if got.Password != nil {
				t.Fatalf("password: got %q, want unset (nil)", *got.Password)
			}

			// Non-write-only fields must still be sourced from the plan.
			if got.Username == nil || *got.Username != "administrator@vsphere.local" {
				t.Fatalf("username: got %v, want administrator@vsphere.local", got.Username)
			}
			if got.ApiUrl != "https://vcenter.example.local" {
				t.Fatalf("api_url: got %q, want https://vcenter.example.local", got.ApiUrl)
			}
		})
	}
}
