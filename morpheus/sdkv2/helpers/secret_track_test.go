// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

// The hashes below are opaque strings as far as the provider is concerned: it
// never computes them, it only compares what the appliance reports with what
// it recorded. Any two distinct values would do.
const (
	configured = "the-configured-secret"
	baseline   = "hash-the-appliance-reported-when-terraform-set-it"
	changed    = "hash-the-appliance-reports-after-somebody-changed-it"
)

func trackedResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"password":      {Type: schema.TypeString, Optional: true, Sensitive: true},
			"password_hash": helpers.SecretHashSchema("password"),
		},
		CustomizeDiff: helpers.SecretHashUnknownWhenSent("password", "password_hash"),
	}
}

func TestUnitSecretHashSchema(t *testing.T) {
	t.Parallel()

	s := helpers.SecretHashSchema("password")
	if !s.Computed || s.Optional || s.Required {
		t.Fatalf("the record must be computed only; got computed=%v optional=%v required=%v", s.Computed, s.Optional, s.Required)
	}

	if !s.Sensitive {
		t.Fatal("the record of a secret's hash must be sensitive")
	}
}

func TestUnitTrackSecret(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		state        map[string]string // password, password_hash before the read
		hash         string            // what the appliance reports now
		sent         bool
		wantPassword string
		wantRecord   string
	}{
		{
			name:         "sent: the hash read back becomes the baseline and the value stays",
			state:        map[string]string{"password": configured},
			hash:         baseline,
			sent:         true,
			wantPassword: configured,
			wantRecord:   baseline,
		},
		{
			name:         "sent again after a drift: the new hash replaces the baseline",
			state:        map[string]string{"password": configured, "password_hash": baseline},
			hash:         changed,
			sent:         true,
			wantPassword: configured,
			wantRecord:   changed,
		},
		{
			name:         "refresh, appliance unchanged: nothing moves",
			state:        map[string]string{"password": configured, "password_hash": baseline},
			hash:         baseline,
			wantPassword: configured,
			wantRecord:   baseline,
		},
		{
			name:         "refresh, changed on the appliance: the hash stands in for the value, the baseline is kept",
			state:        map[string]string{"password": configured, "password_hash": baseline},
			hash:         changed,
			wantPassword: changed,
			wantRecord:   baseline,
		},
		{
			name:         "refresh with no baseline (imported, or state from before the record): the hash stands in",
			state:        map[string]string{"password": configured},
			hash:         baseline,
			wantPassword: baseline,
			wantRecord:   "",
		},
		{
			name:         "the appliance reports no hash: no information, nothing moves",
			state:        map[string]string{"password": configured, "password_hash": baseline},
			hash:         "",
			wantPassword: configured,
			wantRecord:   baseline,
		},
		{
			name:         "sent, but the appliance reports no hash: no baseline is recorded",
			state:        map[string]string{"password": configured},
			hash:         "",
			sent:         true,
			wantPassword: configured,
			wantRecord:   "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			d := trackedResource().TestResourceData()
			for k, v := range c.state {
				if err := d.Set(k, v); err != nil {
					t.Fatal(err)
				}
			}

			if err := helpers.TrackSecret(d, "password", "password_hash", c.hash, c.sent); err != nil {
				t.Fatal(err)
			}

			if got := d.Get("password"); got != c.wantPassword {
				t.Fatalf("password = %q, want %q", got, c.wantPassword)
			}

			if got := d.Get("password_hash"); got != c.wantRecord {
				t.Fatalf("password_hash = %q, want %q", got, c.wantRecord)
			}
		})
	}
}

// After a drift the state holds the appliance's hash, which the configuration
// does not equal: SDKv2 plans an update with no DiffSuppressFunc involved, and
// Update's HasChange is what tells the Read that follows to re-record.
func TestUnitSecretSentOnNewResource(t *testing.T) {
	t.Parallel()

	d := trackedResource().TestResourceData()
	if helpers.SecretSent(d, "password") {
		t.Fatal("a plain refresh has sent nothing")
	}

	d.MarkNewResource()
	if !helpers.SecretSent(d, "password") {
		t.Fatal("a create has sent the value")
	}
}

// The record is known after apply in exactly the plans that send the secret:
// creation, a new configured value, and the configured value going back after
// a drift (where the state holds the appliance's hash). A plan that leaves the
// secret alone leaves the record alone.
func TestUnitSecretHashUnknownWhenSent(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		state       map[string]string // nil: the resource does not exist yet
		config      map[string]any
		wantUnknown bool
	}{
		{
			name:        "create",
			config:      map[string]any{"password": configured},
			wantUnknown: true,
		},
		{
			name:        "new configured value",
			state:       map[string]string{"password": configured, "password_hash": baseline},
			config:      map[string]any{"password": "another-secret"},
			wantUnknown: true,
		},
		{
			name:        "configured value put back after a drift",
			state:       map[string]string{"password": changed, "password_hash": baseline},
			config:      map[string]any{"password": configured},
			wantUnknown: true,
		},
		{
			name:        "unchanged",
			state:       map[string]string{"password": configured, "password_hash": baseline},
			config:      map[string]any{"password": configured},
			wantUnknown: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			var state *terraform.InstanceState
			if c.state != nil {
				state = &terraform.InstanceState{ID: "1", Attributes: c.state}
			}
			diff, err := trackedResource().Diff(context.Background(), state, terraform.NewResourceConfigRaw(c.config), nil)
			if err != nil {
				t.Fatal(err)
			}
			var unknown bool
			if diff != nil {
				if a, ok := diff.Attributes["password_hash"]; ok {
					unknown = a.NewComputed
				}
			}
			if unknown != c.wantUnknown {
				t.Fatalf("password_hash known after apply = %v, want %v (diff %v)", unknown, c.wantUnknown, diff)
			}
		})
	}
}
