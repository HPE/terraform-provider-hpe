// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsource

import (
	"fmt"
	"sort"
	"strings"
)

// Known reports whether source is in the catalogue harvested from Morpheus
// KnownSourcesRelease. A source can be absent and still work — plugins add
// sources at runtime, and newer appliances add more — which is why callers
// warn on an unknown name rather than refuse it.
func Known(source string) bool {
	_, ok := knownSources[source]

	return ok
}

// ServiceTypesFor returns the service types under which source is defined, or
// nil when it is unknown. An entry of "" means the source is reachable with no
// service_type.
func ServiceTypesFor(source string) []string {
	types, ok := knownSources[source]
	if !ok {
		return nil
	}

	out := make([]string, len(types))
	copy(out, types)
	sort.Strings(out)

	return out
}

// Advice returns a warning about the source name and service type as
// configured, or "" when the catalogue has nothing to add.
//
// Two situations are worth a word before the call is made:
//
//   - the source is not in the catalogue. It may still exist (plugin, newer
//     appliance), so this is advisory — but if the call then fails with the
//     generic error, the name is the first thing to check.
//   - the source is defined only under a service type, and the configuration
//     does not set that service type. Without it Morpheus looks on the base
//     service, does not find the source, and fails with the same generic error
//     it gives an unknown name.
func Advice(source, serviceType string) string {
	types := ServiceTypesFor(source)
	if types == nil {
		return fmt.Sprintf("Option source %q is not in the provider's catalogue of sources known for "+
			"Morpheus %s. It may be provided by a plugin or by a newer appliance; if the call fails "+
			"with a generic server error, check the name first.", source, KnownSourcesRelease)
	}

	for _, t := range types {
		if t == serviceType {
			return ""
		}
	}

	nonEmpty := make([]string, 0, len(types))

	for _, t := range types {
		if t != "" {
			nonEmpty = append(nonEmpty, fmt.Sprintf("%q", t))
		}
	}

	if serviceType == "" {
		return fmt.Sprintf("Option source %q is defined for service_type %s and is not reachable without it: "+
			"Morpheus looks on its base service, does not find the source, and returns a generic server "+
			"error. Set service_type.", source, strings.Join(nonEmpty, " or "))
	}

	return fmt.Sprintf("Option source %q is not defined for service_type %q in the catalogue for Morpheus %s "+
		"(it is defined for %s). If the call fails with a generic server error, check service_type.",
		source, serviceType, KnownSourcesRelease, describeTypes(types))
}

func describeTypes(types []string) string {
	parts := make([]string, 0, len(types))

	for _, t := range types {
		if t == "" {
			parts = append(parts, "no service_type")
		} else {
			parts = append(parts, fmt.Sprintf("%q", t))
		}
	}

	return strings.Join(parts, " or ")
}
