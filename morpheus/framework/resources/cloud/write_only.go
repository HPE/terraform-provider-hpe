// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the nested write-only attributes
// config_vmware.password and config_azure.client_secret: see wodigest. A change
// marks write_only_revision unknown, which plans an in-place update in which
// Update re-sends the configured value. The deprecated *_version markers keep
// working for the whole 2.2 line. ModifyPlan is in modify_plan.go.
//
// APIHash adds appliance-side drift detection: Morpheus renders a deterministic
// hash of each stored secret in the cloud's config map (passwordHash for the
// VMware password, clientSecretHash for the Azure client secret). The provider
// records the hash read back after it last sent the value and compares it on
// every later read; an inequality means the secret was changed out of band and
// the configured value is re-sent. The hash is never computed or interpreted.
// The import digest does not seed the appliance baseline; detection starts
// after the first apply that sends the value.
var wo = wodigest.Tracker{
	Resource: "cloud",
	Attributes: map[string]string{ //nolint:gosec // attribute names, not credentials
		"config_vmware.password":     "config_vmware.password_version",
		"config_azure.client_secret": "config_azure.client_secret_version",
	},
	APIHash: map[string]wodigest.HashExtractor{ //nolint:gosec // attribute names, not credentials
		"config_vmware.password":     zoneConfigHashExtractor("passwordHash"),
		"config_azure.client_secret": zoneConfigHashExtractor("clientSecretHash"),
	},
	OnDrift: wodigest.ResendPolicy,
}

// zoneConfigHashExtractor returns a HashExtractor that reads one hash key from
// a cloud's config map. The read response the provider passes is a
// *sdk.GetClouds200ResponseZone; anything else yields present=false.
func zoneConfigHashExtractor(key string) wodigest.HashExtractor {
	return func(resp any) (string, bool) {
		z, ok := resp.(*sdk.GetClouds200ResponseZone)
		if !ok {
			return "", false
		}

		return zoneConfigHash(z, key)
	}
}

// zoneConfigHash reads a single key from a cloud's config map. The config is an
// anyOf whose variants are typed structs, so the hash keys -- which are not
// typed fields -- are reached by marshalling the whole config to JSON (the
// typed variants flatten their AdditionalProperties, where the keys live) and
// reading the key from the resulting map. This is robust to which variant
// matched.
//
// The value is interpreted per the HashExtractor binding rule:
//
//   - no config at all, or the key absent -> ("", false): Morpheus renders a
//     hash key only for a secret the cloud type holds, so a cloud with no such
//     secret simply never carries the key -- no information.
//   - key present and null -> ("", true): a real, comparable value; the secret
//     is unset on the appliance.
//   - key present and a string -> (value, true): the hash.
//   - key present but a non-null, non-string value -> ("", false): the field is
//     not a hash the provider can compare; treat it as no information and log
//     once so the anomaly is visible, rather than coercing it to "" and
//     recording a spurious baseline or drift.
func zoneConfigHash(z *sdk.GetClouds200ResponseZone, key string) (string, bool) {
	if z == nil || z.Config == nil {
		return "", false
	}

	raw, err := json.Marshal(z.Config)
	if err != nil {
		return "", false
	}

	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return "", false
	}

	v, ok := m[key]
	if !ok {
		return "", false
	}

	if v == nil {
		return "", true
	}

	s, ok := v.(string)
	if !ok {
		tflog.Warn(context.Background(),
			"cloud: appliance hash key is not a string; drift detection skipped for it",
			map[string]any{"key": key})

		return "", false
	}

	return s, true
}

// zoneForDrift re-reads the cloud's zone for the drift baseline. Drift
// detection is best-effort: a failed or empty read yields nil, and
// StoreObserved treats a nil response as "nothing to baseline" rather than an
// error, so the apply is never blocked by it. Read reuses the zone it already
// fetched instead of calling this; it is used by Create (post read-back) and
// Update, where the read-back response is a different type.
func zoneForDrift(ctx context.Context, id int64, client *sdk.APIClient) *sdk.GetClouds200ResponseZone {
	c, _, err := client.CloudsAPI.GetClouds(ctx, id).Execute()
	if err != nil || c == nil {
		return nil
	}

	return c.Zone
}
