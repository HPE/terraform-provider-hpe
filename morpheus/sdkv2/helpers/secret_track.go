// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A secret the appliance reports back only as a hash -- an integration's
// password or token -- is tracked the way the framework resources track
// write-only values (morpheus/utils/wodigest): the provider never computes or
// interprets the hash. It records the one read back right after it last sent
// the value, and compares the hash of every later read with that record for
// string equality; an inequality means the secret changed on the appliance
// after Terraform set it. SDKv2 resources have no private state, so the record
// lives in a Computed, Sensitive <attribute>_hash beside the attribute, which
// SecretHashSchema declares.

// SecretHashSchema is the schema of the attribute that records the appliance's
// hash of the secret named by attribute, for TrackSecret.
func SecretHashSchema(attribute string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		// A hash of a secret is still about the secret; it is kept out of
		// plan output.
		Sensitive: true,
		Description: "The appliance's hash of " + attribute + " as read back when Terraform last set it. " +
			"Each read compares the hash the appliance reports with it; a difference means the value was " +
			"changed outside Terraform, and the next plan sets the configured value again.",
	}
}

// TrackSecret updates the attribute named by key, and its record hashKey, from
// the hash the appliance reports for the secret now.
//
// sent says whether the value was sent in this apply -- the resource is being
// created, or the update that is finishing changed it (SecretSent): the hash
// read back is then the baseline, recorded in hashKey, and the configured
// value stays in state. Otherwise the hash is compared with the baseline.
// Equal means the appliance still holds what Terraform set, and the state
// value stays. Different -- changed out of band, or no baseline at all because
// the resource was imported or its state predates the record -- puts the
// appliance's hash in the attribute's place in state: the configured value
// does not equal it, so the next plan sends the configured value, and the
// apply records the new baseline. An empty hash is no information (the
// appliance holds no such secret, or does not report one) and changes
// nothing.
func TrackSecret(d *schema.ResourceData, key, hashKey, hash string, sent bool) error {
	if hash == "" {
		return nil
	}

	if sent {
		return d.Set(hashKey, hash)
	}

	if recorded, _ := d.Get(hashKey).(string); recorded == hash {
		return nil
	}

	return d.Set(key, hash)
}

// SecretSent reports, for the Read a Create or Update runs when it is done,
// whether the value named by key was sent in this apply: the resource is new,
// or the update changed the attribute (which is also the case after
// TrackSecret put the appliance's hash in its place, so a drifted value is
// re-sent and re-recorded). On a plain refresh there is no plan to differ
// from, so it is false.
func SecretSent(d *schema.ResourceData, key string) bool {
	return d.IsNewResource() || d.HasChange(key)
}

// SecretHashUnknownWhenSent is the CustomizeDiff that goes with TrackSecret:
// a plan that changes the attribute named by key -- a new value, or the
// configured value put back after a drift -- is going to send it, and the
// apply will record the hash the appliance then reports, so the plan shows
// hashKey as known after apply rather than as its current value. Compose it
// with the resource's other CustomizeDiff functions with customdiff.All.
func SecretHashUnknownWhenSent(key, hashKey string) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(hashKey, func(_ context.Context, d *schema.ResourceDiff, _ any) bool {
		return d.HasChange(key)
	})
}
