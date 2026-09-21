// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsource

import "fmt"

// contextHint describes a source that answers with an empty list, rather than an
// error, when the parameters that give it a provisioning context are absent.
type contextHint struct {
	// anyOf are parameter keys of which at least one must be present; when none
	// is, the hint applies.
	anyOf []string
	// advice names the data-source attributes that supply them.
	advice string
}

// contextHints is a curated table, not a validator. It holds only sources whose
// empty-list-on-missing-context behaviour has been observed on an appliance; an
// unknown source gets no hint. It is expected to stay short and to grow from
// incidents.
var contextHints = map[string]contextHint{
	// The provisioning wizard's resource-pool list. Without a layout or
	// provision type it cannot tell which pools are provisionable, and returns
	// nothing. (Observed: zoneId alone -> 0 rows; with layoutId -> the cluster pools.)
	"zonePools": {
		anyOf:  []string{"layoutId", "provisionTypeCode", "provisionType", "provisionTypeId"},
		advice: "set layout_id or provision_type_code in addition to cloud_id",
	},
	// The provisioning wizard's network list. supportsNetworkSelection is read
	// from the provision type, which comes from the layout; without one the
	// networks list is empty. (Observed: zoneId alone -> networks: 0; with
	// layoutId -> the pool's networks.)
	"zoneNetworkOptions": {
		anyOf:  []string{"layoutId", "provisionType", "provisionTypeId", "provisionTypeCode"},
		advice: "set layout_id or provision_type_code in addition to cloud_id (and pool_id to scope to a resource pool)",
	},
}

// MissingContextHint returns advice when source returned nothing and belongs to
// the family that returns an empty list, not an error, without its context
// parameters — and those parameters were indeed absent. Otherwise "".
func MissingContextHint(source string, params map[string]string) string {
	h, ok := contextHints[source]
	if !ok {
		return ""
	}

	for _, k := range h.anyOf {
		if v, present := params[k]; present && v != "" {
			return ""
		}
	}

	return fmt.Sprintf("%s requires a provisioning context: %s. Without one the server returns an empty "+
		"list rather than an error.", source, h.advice)
}
