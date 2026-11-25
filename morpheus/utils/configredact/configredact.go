// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package configredact withholds credentials from the free-form configuration
// maps Morpheus returns, before they are written to Terraform state.
//
// Morpheus redacts the well-known top-level secrets of an object
// (servicePassword, serviceToken, serviceAccess …) but returns the object's
// config map largely as stored, and some types keep credentials in it — an
// External Kubernetes cluster's config carries apiToken (a cluster-admin JWT)
// and serviceAccess (the kubeconfig holding it). The appliance's own redaction
// list for a cluster's config knows only kubeToken. A provider read copying
// that map into state would publish live credentials in terraform.tfstate and
// in plan output, so the keys are removed here — removed rather than masked,
// since a masked value would still occupy the attribute and would never
// match the configuration.
package configredact

import (
	"regexp"
	"sort"
	"strings"
)

// clusterExact are config keys of a cluster known to hold credentials or
// material that embeds them, whatever their name suggests. Compared
// case-insensitively, so a differently cased spelling does not slip through.
var clusterExact = lowerSet(
	"apiToken",
	"serviceAccess",
	"serviceToken",
	"kubeToken",
	"kubeConfig",
	"joinCommand",
)

func lowerSet(keys ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}

	return m
}

// credentialName matches key names that end in a word for a credential, so a
// key this package has not met is still withheld when its name says what it
// holds. Case-insensitive; the suffix form keeps names such as "tokenExpiry"
// or "passwordPolicy" out of it.
var credentialName = regexp.MustCompile(`(?i)(token|password|passwd|secret|passphrase|privatekey|kubeconfig|serviceaccess)$`)

// Cluster returns a copy of a cluster's config with credential keys removed,
// at any depth, and the names of the keys removed (top level and nested,
// sorted, without duplicates). A nil map returns nil. The config is what the
// SDK decoded from JSON, so nested values are map[string]any and []any; the
// typed shapes a caller might build by hand are handled as well.
func Cluster(config map[string]any) (map[string]any, []string) {
	if config == nil {
		return nil, nil
	}

	seen := map[string]struct{}{}
	out := redactMap(config, clusterExact, seen)

	withheld := make([]string, 0, len(seen))
	for k := range seen {
		withheld = append(withheld, k)
	}

	sort.Strings(withheld)

	return out, withheld
}

func isCredentialKey(key string, exact map[string]struct{}) bool {
	if _, ok := exact[strings.ToLower(key)]; ok {
		return true
	}

	return credentialName.MatchString(key)
}

func redactMap(m map[string]any, exact map[string]struct{}, seen map[string]struct{}) map[string]any {
	out := make(map[string]any, len(m))

	for k, v := range m {
		if isCredentialKey(k, exact) {
			seen[k] = struct{}{}

			continue
		}

		out[k] = redactValue(v, exact, seen)
	}

	return out
}

func redactValue(v any, exact map[string]struct{}, seen map[string]struct{}) any {
	switch t := v.(type) {
	case map[string]any:
		return redactMap(t, exact, seen)
	case map[string]string:
		out := make(map[string]string, len(t))
		for k, e := range t {
			if isCredentialKey(k, exact) {
				seen[k] = struct{}{}

				continue
			}

			out[k] = e
		}

		return out
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = redactValue(e, exact, seen)
		}

		return out
	case []map[string]any:
		out := make([]map[string]any, len(t))
		for i, e := range t {
			out[i] = redactMap(e, exact, seen)
		}

		return out
	default:
		return v
	}
}
