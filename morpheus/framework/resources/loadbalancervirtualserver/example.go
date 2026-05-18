// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_virtual_server/example_nsxt_full.tf example_nsxt_full.tf.tmpl LoadBalancerId 1 VipName "example-nsxt-vip-ssl-client" Description "Example NSX-T virtual server" VipAddress "10.0.0.5" VipPort 443 VipProtocol "http" PoolId 11 SslCert 12 SslServerCert 0 ApplicationProfile 13 Persistence "COOKIE" PersistenceProfile 16 SslClientProfile 19 SslServerProfile 0
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_virtual_server/example_nsxt_minimal.tf example_nsxt_minimal.tf.tmpl LoadBalancerId 1 VipName "example-nsxt-vip" Description "Example NSX-T virtual server" VipAddress "10.0.0.4" VipPort 443 VipProtocol "http" PoolId 11 ApplicationProfile 13

func RenderLoadBalancerVirtualServerNsxtFullConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":     "1",
		"VipName":            "example-nsxt-vip-ssl-client",
		"Description":        "Example NSX-T virtual server",
		"VipAddress":         "10.0.0.5",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "11",
		"SslCert":            "12",
		"SslServerCert":      "0",
		"ApplicationProfile": "13",
		"Persistence":        "COOKIE",
		"PersistenceProfile": "16",
		"SslClientProfile":   "19",
		"SslServerProfile":   "0",
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
	templatePath := filepath.Join(dir, "example_nsxt_full.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderLoadBalancerVirtualServerNsxtMinimalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":     "1",
		"VipName":            "example-nsxt-vip",
		"Description":        "Example NSX-T virtual server",
		"VipAddress":         "10.0.0.4",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "11",
		"ApplicationProfile": "13",
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
	templatePath := filepath.Join(dir, "example_nsxt_minimal.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
