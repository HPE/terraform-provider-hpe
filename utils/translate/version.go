// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

// ResolveForVersion returns a new CompiledConfig with version-specific
// overrides merged in. If ver is empty or no version overrides match,
// the base config is returned unchanged.
func (cc *CompiledConfig) ResolveForVersion(ver string) *CompiledConfig {
	if ver == "" || len(cc.raw.Versions) == 0 {
		return cc
	}

	morpheusVersion, err := version.NewVersion(strings.TrimPrefix(ver, "v"))
	if err != nil {
		// If version string is unparseable, return base config
		return cc
	}

	// Collect matching version overrides
	var extraMoves []map[string]string
	var extraRemoves []string

	for _, vo := range cc.raw.Versions {
		constraint, err := version.NewConstraint(vo.Constraint)
		if err != nil {
			continue
		}

		if constraint.Check(morpheusVersion) {
			extraMoves = append(extraMoves, vo.Moves...)
			extraRemoves = append(extraRemoves, vo.Removes...)
		}
	}

	if len(extraMoves) == 0 && len(extraRemoves) == 0 {
		return cc
	}

	// Build a new config with the merged moves/removes
	mergedCfg := &ResourceConfig{
		Templates:     cc.raw.Templates,
		Moves:         append(append([]map[string]string{}, cc.raw.Moves...), extraMoves...),
		Removes:       append(append([]string{}, cc.raw.Removes...), extraRemoves...),
		Overrides:     cc.raw.Overrides,
		Versions:      nil, // Don't recurse
		Paths:         cc.raw.Paths,
		Envelope:      cc.raw.Envelope,
		Discriminator: cc.raw.Discriminator,
	}

	return Compile(mergedCfg)
}

// MustParseConstraint parses a version constraint or panics.
func MustParseConstraint(s string) version.Constraints {
	c, err := version.NewConstraint(s)
	if err != nil {
		panic(fmt.Sprintf("invalid version constraint: %s", s))
	}

	return c
}
