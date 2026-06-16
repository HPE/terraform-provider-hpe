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

// renderNsxtTier1Prereq renders a per-test NSX-T tier-1 gateway connected to an
// existing tier-0 (router id 28) and a data source exposing the tier-1's
// provider_id for use as a load balancer's tier1_gateway (data source
// hpe_morpheus_network_router.lb_tier1).
//
// NSX-T allows only one load balancer service per tier-1, so each test must use
// its own tier-1 (sharing one would collide under parallel runs). The tier-1 is
// connected to the pre-provisioned tier-0 (whose provider_id/path we read via a
// data source) and given an edge cluster, both required for an LB service to
// deploy on it.
//
// QA verify: tier-0 router id 28 is a realized NSX-T tier-0 on integration 5;
// edge_cluster is the NSX-T edge cluster external id (display name
// "qa-edge-cluster-01").
func renderNsxtTier1Prereq(t *testing.T, name string) (string, error) {
	t.Helper()

	return `
data "hpe_morpheus_network_router" "lb_tier0" {
  id = 28
}

resource "hpe_morpheus_network_router" "lb_tier1" {
  name                   = "` + name + `-tier1"
  group_id               = 3
  network_integration_id = 5

  config_nsxt_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
    edge_cluster       = "3de5f8d0-4f8a-433b-95ed-91020c948084"
    fail_over          = "NON_PREEMPTIVE"
    tier0_gateway      = data.hpe_morpheus_network_router.lb_tier0.provider_id
  }
}

data "hpe_morpheus_network_router" "lb_tier1" {
  id = hpe_morpheus_network_router.lb_tier1.id
}
`, nil
}
