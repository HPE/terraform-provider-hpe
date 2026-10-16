// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
)

// TestUnitExampleDataSourceAddressesMatchTemplate pins the Example*Data
// constants to the data sources example.tf.tmpl declares, so a rename in the
// template fails here rather than as an acceptance failure in a package that
// builds on the rendered configuration.
func TestUnitExampleDataSourceAddressesMatchTemplate(t *testing.T) {
	t.Parallel()

	tmpl, err := os.ReadFile("example.tf.tmpl")
	if err != nil {
		t.Fatal(err)
	}

	declared := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^(data|resource) "([a-z_]+)" "([a-z_]+)"`).FindAllStringSubmatch(string(tmpl), -1) {
		addr := m[2] + "." + m[3]
		if m[1] == "data" {
			addr = "data." + addr
		}

		declared[addr] = true
	}

	for _, want := range []string{
		instance.ExampleCloudData,
		instance.ExampleClusterData,
		instance.ExampleLayoutData,
		instance.ExampleNetworkData,
		instance.ExampleDatastoreData,
		instance.ExamplePlanData,
		instance.ExampleImageData,
		instance.ExampleGroupData,
		instance.ExampleInstance,
	} {
		if !declared[want] {
			t.Errorf("%s is not declared in example.tf.tmpl (declared: %s)", want, strings.Join(keys(declared), ", "))
		}
	}
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	return out
}
