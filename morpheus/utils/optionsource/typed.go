// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsource

import (
	"fmt"
	"sort"
	"strings"
)

// Typed holds the data sources' typed convenience attributes. Each maps to the
// parameter key(s) the sources read, so a practitioner writes cloud_id, the
// provider's name, rather than remembering that the API calls it zoneId — and,
// for the attributes Morpheus reads under more than one key, does not have to
// know that either.
//
// A nil pointer means the attribute was not set.
type Typed struct {
	CloudID           *int64
	Phrase            *string
	ParentID          *int64
	CredentialID      *int64
	GroupID           *int64
	LayoutID          *int64
	ProvisionTypeCode *string
	ProvisionTypeID   *int64
	PoolID            *string
	PlanID            *int64
	ServerGroupID     *int64
	ProviderID        *int64
}

// typedKeys maps each typed attribute to the parameter keys it sets.
//
// The rule for which keys get a typed attribute: a key that is read as a
// scalar id or string — something an API caller can supply — by roughly eight
// or more sources, checked against the code that consumes it rather than
// against how often the name appears. Around two hundred distinct keys are
// read across the catalogue; the rest are reached through the parameters map.
// Keys the sources read as a hydrated server-side object (domain, parent,
// report, router, network, config, ...) are never typed: the UI passes those
// internally and an API caller cannot.
//
// Attribute names follow the provider's convention, not the API's: the API
// says zone where the provider says cloud, and the mapping here is where the
// translation happens. Where sources read the same thing under more than one
// name, the attribute is sent under every name:
//
//   - cloud_id: zoneId, and cloudId, which newer sources read for the same id.
//   - group_id: groupId and siteId.
//   - parent_id: parentId and parentid. Never the bare "parent", which is one
//     of the object-valued keys.
//   - provision_type_code: provisionTypeCode and provisionType. The id form is
//     a separate attribute, provision_type_id → provisionTypeId, because it is
//     a different value.
//   - credential_id: credential.type, which the server's nested parameter
//     binding turns into the credential map its loader expects.
//   - pool_id accepts either "pool-<id>" or a bare id; the server strips the
//     prefix itself.
var typedKeys = map[string][]string{
	"cloud_id":            {"zoneId", "cloudId"},
	"phrase":              {"phrase"},
	"parent_id":           {"parentId", "parentid"},
	"credential_id":       {"credential.type"},
	"group_id":            {"groupId", "siteId"},
	"layout_id":           {"layoutId"},
	"provision_type_code": {"provisionTypeCode", "provisionType"},
	"provision_type_id":   {"provisionTypeId"},
	"pool_id":             {"poolId"},
	"plan_id":             {"planId"},
	"server_group_id":     {"serverGroupId"},
	"provider_id":         {"providerId"},
}

// TypedKeys returns the parameter keys a typed attribute sets, for documentation.
func TypedKeys(attribute string) []string {
	return typedKeys[attribute]
}

// serverSetKeys are the parameters the server sets for every option-source
// call from the authenticated user, discarding any value the caller supplied
// (OptionSourceService.getOptionSourceParams, present since at least 8.0).
// accountId is the one a practitioner might reasonably send, expecting to read
// a source as another tenant; there is no way to do that on this endpoint
// other than to authenticate as a user of that tenant. The rest are internal.
var serverSetKeys = []string{
	"accountId", "user", "currentUser", "currentUserEmail", "currentUsername",
	"jsonPayload", "serviceType", "namespace", "requestSourceType",
}

// ServerSetParameters returns, in a stable order, the keys of params that the
// server will replace with its own values, so the caller can warn that they
// have no effect. The tenant is the case that matters: accountId is always
// the caller's tenant, whatever is sent.
func ServerSetParameters(params map[string]string) []string {
	var out []string

	for _, k := range serverSetKeys {
		if _, ok := params[k]; ok {
			out = append(out, k)
		}
	}

	return out
}

// ServerSetAdvice is the warning text for one key ServerSetParameters returned.
func ServerSetAdvice(key string) string {
	if key == "accountId" {
		return "Morpheus sets accountId to the authenticated user's tenant on every option-source call and " +
			"discards the value supplied, so this source is read for the caller's tenant. To read it as " +
			"another tenant, authenticate as a user of that tenant."
	}

	return fmt.Sprintf("Morpheus sets %s itself on every option-source call and discards the value supplied; "+
		"it has no effect here.", key)
}

// MergeParameters combines the free-form parameters map with the typed
// attributes into the map the request sends.
//
// A key supplied both ways is an error rather than a silent precedence rule:
// the two spellings are equivalent, so a conflict is always a mistake the
// author should see.
func MergeParameters(parameters map[string]string, typed Typed) (map[string]string, error) {
	out := make(map[string]string, len(parameters)+12)
	for k, v := range parameters {
		out[k] = v
	}

	set := func(attribute, value string) error {
		for _, key := range typedKeys[attribute] {
			if existing, dup := out[key]; dup {
				return fmt.Errorf("parameter %q is set both by %s and in parameters (%q); use one",
					key, attribute, existing)
			}

			out[key] = value
		}

		return nil
	}

	type entry struct {
		attribute string
		value     string
		present   bool
	}

	entries := []entry{
		{"cloud_id", int64String(typed.CloudID), typed.CloudID != nil},
		{"phrase", deref(typed.Phrase), typed.Phrase != nil},
		{"parent_id", int64String(typed.ParentID), typed.ParentID != nil},
		{"credential_id", int64String(typed.CredentialID), typed.CredentialID != nil},
		{"group_id", int64String(typed.GroupID), typed.GroupID != nil},
		{"layout_id", int64String(typed.LayoutID), typed.LayoutID != nil},
		{"provision_type_code", deref(typed.ProvisionTypeCode), typed.ProvisionTypeCode != nil},
		{"provision_type_id", int64String(typed.ProvisionTypeID), typed.ProvisionTypeID != nil},
		{"pool_id", deref(typed.PoolID), typed.PoolID != nil},
		{"plan_id", int64String(typed.PlanID), typed.PlanID != nil},
		{"server_group_id", int64String(typed.ServerGroupID), typed.ServerGroupID != nil},
		{"provider_id", int64String(typed.ProviderID), typed.ProviderID != nil},
	}

	for _, e := range entries {
		if !e.present {
			continue
		}

		if err := set(e.attribute, e.value); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// DescribeParameters renders a parameter map for a message, keys sorted.
func DescribeParameters(params map[string]string) string {
	if len(params) == 0 {
		return "(none)"
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	return strings.Join(parts, ", ")
}

func int64String(p *int64) string {
	if p == nil {
		return ""
	}

	return fmt.Sprintf("%d", *p)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}
