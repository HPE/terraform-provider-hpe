// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the four write-only credentials: see wodigest. Their
// configured values are never stored in state; their digests are kept in
// private state, and a change to one marks write_only_revision unknown, which
// plans an in-place update in which Update re-sends the configured values.
//
// APIHash adds detection of a change made on the appliance (UI, API, another
// tool). For the two passwords Morpheus returns, beside the masked value, a
// deterministic hash of the stored secret (smtpPasswordHash,
// proxyPasswordHash); the currency key it returns in clear, so its extractor
// hashes the key itself and private state never holds it. Each hash is
// recorded right after the provider sends a value and only compared for
// equality on later reads -- its algorithm is never assumed. OnDrift re-sends
// the configured value when the stored one changed. The Twilio auth token has
// no extractor: Morpheus returns it masked and with no hash, so only a change
// in the configuration is detected.
//
// The attributes keep the names they had in the SDKv2 resource (no _wo
// suffix), so existing configurations are unchanged, and have no deprecated
// *_version markers; twilio_auth_token, which the SDKv2 resource did not have,
// is named as the spec names it.
var wo = wodigest.Tracker{
	Resource: resourceType,
	Attributes: map[string]string{ //nolint:gosec // attribute names, not credentials
		"smtp_password":             "",
		"proxy_password":            "",
		"currency_provider_api_key": "",
		"twilio_auth_token":         "",
	},
	APIHash: map[string]wodigest.HashExtractor{ //nolint:gosec // attribute names, not credentials
		"smtp_password":             smtpPasswordHash,
		"proxy_password":            proxyPasswordHash,
		"currency_provider_api_key": currencyKeyDigest,
	},
	OnDrift: wodigest.ResendPolicy,
}

// smtpPasswordHash is the appliance's hash of the stored SMTP password from a
// settings read. present=false when the response is not a settings object or
// the API omitted the field (no information); a field rendered null -- no
// password stored -- is ("", true).
func smtpPasswordHash(resp any) (string, bool) {
	s, ok := resp.(*applianceSettings)
	if !ok || s == nil {
		return "", false
	}

	return wodigest.NullableStringHash(s.SmtpPasswordHash)
}

// proxyPasswordHash is smtpPasswordHash for the proxy password.
func proxyPasswordHash(resp any) (string, bool) {
	s, ok := resp.(*applianceSettings)
	if !ok || s == nil {
		return "", false
	}

	return wodigest.NullableStringHash(s.ProxyPasswordHash)
}

// currencyKeyDigest is the sha256 hex digest of the currency provider key the
// appliance holds, which the API returns in clear: the digest, not the key, is
// what drift detection records in private state. A key the API renders null
// or empty is no key, ("", true); an omitted field is no information.
func currencyKeyDigest(resp any) (string, bool) {
	s, ok := resp.(*applianceSettings)
	if !ok || s == nil || !s.CurrencyKey.IsSet() {
		return "", false
	}

	key := s.CurrencyKey.Get()
	if key == nil || *key == "" {
		return "", true
	}

	sum := sha256.Sum256([]byte(*key))

	return hex.EncodeToString(sum[:]), true
}

// ModifyPlan plans the update that re-sends the credentials when one changed
// in the configuration or on the appliance, or when the resource has no record
// of them at all, as after an upgrade from the SDKv2 resource; see
// write_only_revision. An import is not that case: it records the digests its
// ID carries, or a sentinel that takes the configured credentials as already
// applied (see import.go), so the plan after an import re-sends nothing.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
