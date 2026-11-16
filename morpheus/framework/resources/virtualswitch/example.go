// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_virtual_switch/resource.tf example.tf.tmpl ClusterId "1" Name "vsw-hvm" BondMode "active-backup" Mtu "9000" Nic1 "eth2" Nic2 "eth3" NfsVlan "100" NfsNetmask "255.255.255.0" NfsIpExpr "local.nfs_ips[host.server_name]"

// RenderVirtualSwitchExample renders the two-host virtual switch example with
// the given overrides; used by the example generator and tests.
func RenderVirtualSwitchExample(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ClusterId":  "1",
		"Name":       "vsw-hvm",
		"BondMode":   "active-backup",
		"Mtu":        "9000",
		"Nic1":       "eth2",
		"Nic2":       "eth3",
		"NfsVlan":    "100",
		"NfsNetmask": "255.255.255.0",
		// The test cannot know host names, so it derives each address from
		// the host's position among the sorted host names (a server id can
		// exceed a /24); the documentation example uses an explicit map by name.
		"NfsIpExpr": `cidrhost("10.20.0.0/24", 10 + index(local.host_names, host.server_name))`,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

// TestNamePrefix is the name prefix of virtual switches created by the
// acceptance tests. Switch names are limited to 16 characters, so the
// suite-wide sweep prefix cannot be used. The sweeper (build tag "sweep")
// holds its own copy of the literal because it cannot import this package
// without an import cycle; TestUnitSweepPrefixMatches keeps the two in step.
const TestNamePrefix = "tfacc"
