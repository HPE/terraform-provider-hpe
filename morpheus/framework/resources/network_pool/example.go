// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_pool

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_network_pool/example.tf example.tf.tmpl Name "App Pool" TypeId "1" SubnetAddress "10.0.1.0" Netmask "255.255.255.0" Gateway "10.0.1.1" DnsDomain "example.com"

func RenderNetworkPoolConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":          "App Pool",
		"TypeId":        "1",
		"SubnetAddress": "10.0.1.0",
		"Netmask":       "255.255.255.0",
		"Gateway":       "10.0.1.1",
		"DnsDomain":     "example.com",
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

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
