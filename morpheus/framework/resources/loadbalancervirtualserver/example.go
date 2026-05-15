// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_virtual_server/example.tf example.tf.tmpl LoadBalancerId 1 VipName "example-vip" Description "Example virtual server" VipAddress "10.0.0.1" VipPort 80 VipProtocol "http"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_virtual_server/example_nsxt.tf example_nsxt.tf.tmpl LoadBalancerId 1 VipName "example-nsxt-vip" Description "Example NSX-T virtual server" VipAddress "10.0.0.2" VipPort 443 VipProtocol "http" PoolId 42 SslCert 12 SslServerCert 0 ApplicationProfile 85 Persistence "SOURCE_IP" PersistenceProfile 78 SslClientProfile 33 SslServerProfile 0

func RenderLoadBalancerVirtualServerConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId": "1",
		"VipName":        "example-vip",
		"Description":    "Example virtual server",
		"VipAddress":     "10.0.0.1",
		"VipPort":        "80",
		"VipProtocol":    "http",
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

func RenderLoadBalancerVirtualServerNsxtConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":     "1",
		"VipName":            "example-nsxt-vip",
		"Description":        "Example NSX-T virtual server",
		"VipAddress":         "10.0.0.2",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "42",
		"SslCert":            "12",
		"SslServerCert":      "0",
		"ApplicationProfile": "85",
		"Persistence":        "SOURCE_IP",
		"PersistenceProfile": "78",
		"SslClientProfile":   "33",
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
	templatePath := filepath.Join(dir, "example_nsxt.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
