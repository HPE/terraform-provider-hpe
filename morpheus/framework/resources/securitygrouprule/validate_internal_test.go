// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygrouprule

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func withSourceType(m SecurityGroupRuleModel, typ string) SecurityGroupRuleModel {
	if typ == "" {
		m.SourceType = types.StringUnknown()
	} else {
		m.SourceType = types.StringValue(typ)
	}

	return m
}

func withDestinationType(m SecurityGroupRuleModel, typ string) SecurityGroupRuleModel {
	m.DestinationType = types.StringValue(typ)

	return m
}

func TestUnitValidateRequiredByCondition(t *testing.T) {
	t.Parallel()

	s := types.StringValue
	null := types.StringNull()
	unknown := types.StringUnknown()

	cidr := s("10.0.0.0/24")
	any := s("0.0.0.0/0")

	rule := func(proto, src, dst, ports types.String) SecurityGroupRuleModel {
		return SecurityGroupRuleModel{Protocol: proto, Source: src, Destination: dst, PortRange: ports}
	}

	cases := []struct {
		name  string
		m     SecurityGroupRuleModel
		wants []string // attributes expected in errors, in any order
	}{
		{"complete cidr tcp rule", rule(s("tcp"), cidr, any, s("443")), nil},
		{
			"types unset default to cidr: all three missing", rule(s("tcp"), null, null, null),
			[]string{"source", "destination", "port_range"},
		},
		{"icmp needs no port range", rule(s("icmp"), cidr, any, null), nil},
		{"udp needs a port range", rule(s("udp"), cidr, any, null), []string{"port_range"}},
		{"group source needs no cidr", withSourceType(rule(s("tcp"), null, any, s("22")), "group"), nil},
		{"tier destination needs no cidr", withDestinationType(rule(s("tcp"), cidr, null, s("22")), "tier"), nil},
		{"explicit cidr type, source missing", withSourceType(rule(s("tcp"), null, any, s("22")), "cidr"), []string{"source"}},
		{"unknown source is not judged", rule(s("tcp"), unknown, any, s("22")), nil},
		{"unknown source type is not judged", withSourceType(rule(s("tcp"), null, any, s("22")), ""), nil},
		{"unknown protocol is not judged", rule(unknown, cidr, any, null), nil},
		{"unknown port range is not judged", rule(s("tcp"), cidr, any, unknown), nil},
		{"null protocol: the framework reports the required attribute", rule(null, cidr, any, null), nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diags := validateRequiredByCondition(tc.m)

			var got []string
			for _, d := range diags.Errors() {
				// The attribute path is the last segment of the summary's target.
				got = append(got, strings.TrimPrefix(d.Summary(), "Missing "))
			}

			if len(got) != len(tc.wants) {
				t.Fatalf("errors = %v, want %v", got, tc.wants)
			}

			for _, w := range tc.wants {
				found := false
				for _, g := range got {
					if g == w {
						found = true
					}
				}

				if !found {
					t.Fatalf("missing an error on %q; got %v", w, got)
				}
			}
		})
	}
}
