// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_pool/example.tf example.tf.tmpl LoadBalancerId "1" Name "NSX-T Pool" Description "An NSX-T load balancer pool" VipBalance "ROUND_ROBIN" MinActive "1" SnatTranslationType "LBSnatAutoMap" TcpMultiplexing "true" TcpMultiplexingNumber "6"

func RenderLoadBalancerPoolNsxtConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":        "1",
		"Name":                  "NSX-T Pool",
		"Description":           "An NSX-T load balancer pool",
		"VipBalance":            "ROUND_ROBIN",
		"MinActive":             "1",
		"SnatTranslationType":   "LBSnatAutoMap",
		"TcpMultiplexing":       "true",
		"TcpMultiplexingNumber": "6",
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
