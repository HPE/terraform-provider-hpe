// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer/example_haproxy.tf example_haproxy.tf.tmpl Name "example-terraform-haproxy-lb" Description "HAProxy load balancer" CloudName "hvm" GroupName "Zodiac" PlanId "8" Pool "pool-574" Visibility "public"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer/example_haproxy_generic.tf example_haproxy_generic.tf.tmpl Name "example-terraform-haproxy-lb" Description "HAProxy load balancer via generic config" CloudName "hvm" GroupName "Zodiac" TypeCode "haproxyContainer" PlanId "8" Pool "pool-574" Visibility "public"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer/example_nsxt.tf example_nsxt.tf.tmpl Name "example-terraform-nsxt-lb" TypeCode "nsx-t" Visibility "public" AdminState "true" LogLevel "INFO" Size "SMALL" Tier1Gateway "\"tier1-gateway\""

func RenderLoadBalancerHAProxyConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":   "hvm",
		"GroupName":   "Zodiac",
		"Name":        "example-terraform-haproxy-lb",
		"Description": "HAProxy load balancer",
		"PlanId":      "8",
		"Pool":        "pool-574",
		"Visibility":  "public",
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
	templatePath := filepath.Join(dir, "example_haproxy.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderLoadBalancerHAProxyGenericConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":   "hvm",
		"GroupName":   "Zodiac",
		"Name":        "example-terraform-haproxy-lb",
		"Description": "HAProxy load balancer via generic config",
		"TypeCode":    "haproxyContainer",
		"PlanId":      "8",
		"Pool":        "pool-574",
		"Visibility":  "public",
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
	templatePath := filepath.Join(dir, "example_haproxy_generic.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderLoadBalancerNsxtConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "example-terraform-nsxt-lb",
		"TypeCode":   "nsx-t",
		"Visibility": "public",
		"AdminState": "true",
		"LogLevel":   "INFO",
		"Size":       "SMALL",
		// Self-contained tier-1 gateway: the LB's tier1_gateway expects the
		// provider_id of an NSX-T tier-1 gateway, which is exposed by the
		// hpe_morpheus_network_router data source created in the prereq below.
		"Tier1Gateway": "data.hpe_morpheus_network_router.lb_tier1.provider_id",
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

	lb, err := testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
	if err != nil {
		return "", err
	}

	prereq, err := renderNsxtTier1Prereq(t, defaults["Name"])
	if err != nil {
		return "", err
	}

	return prereq + lb, nil
}

// renderNsxtTier1Prereq renders a data source exposing the provider_id of an
// existing, fully-realized NSX-T tier-1 gateway (router id 27) for use as a load
// balancer's tier1_gateway (data source hpe_morpheus_network_router.lb_tier1).
//
// An LB service can only be deployed on a tier-1 that is connected to a tier-0
// (or has a Tier1Interface) and has an associated edge cluster. Rather than
// building that topology per test (and racing NSX-T realization), we reference a
// pre-provisioned tier-1.
//
// QA verify: router id 27 is a realized NSX-T tier-1 (connected to a tier-0, with
// an edge cluster) on integration 5.
func renderNsxtTier1Prereq(t *testing.T, _ string) (string, error) {
	t.Helper()

	return `
data "hpe_morpheus_network_router" "lb_tier1" {
  id = 27
}
`, nil
}
