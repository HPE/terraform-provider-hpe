// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_monitor/example_nsxt.tf example_nsxt.tf.tmpl LoadBalancerId "1" Name "NSX-T HTTP Monitor" Description "An NSX-T HTTP health check monitor" MonitorType "http" MonitorInterval "5" MonitorTimeout "15" MonitorDestination "/" FallCount "3" RiseCount "3" AliasPort "8080" SendData "GET / HTTP/1.1" SendType "GET" SendVersion "HTTP_VERSION_1_1" ReceiveData "" ReceiveCode "200" DataLength "0" MaxRetry "3"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_monitor/example_nsxv.tf example_nsxv.tf.tmpl LoadBalancerId "1" Name "NSX-V HTTP Monitor" Description "An NSX-V HTTP health check monitor" MonitorType "http" MonitorInterval "10" MonitorTimeout "15" MaxRetry "3" SendData "GET / HTTP/1.0" SendType "GET" ReceiveData "" ReceiveCode "200" MonitorDestination "/health" MonitorUsername "" MonitorPassword ""

func RenderLoadBalancerMonitorNsxtConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":     "1",
		"Name":               "NSX-T HTTP Monitor",
		"Description":        "An NSX-T HTTP health check monitor",
		"MonitorType":        "http",
		"MonitorInterval":    "5",
		"MonitorTimeout":     "15",
		"MonitorDestination": "/",
		"FallCount":          "3",
		"RiseCount":          "3",
		"AliasPort":          "8080",
		"SendData":           "GET / HTTP/1.1",
		"SendType":           "GET",
		"SendVersion":        "HTTP_VERSION_1_1",
		"ReceiveData":        "",
		"ReceiveCode":        "200",
		"DataLength":         "0",
		"MaxRetry":           "3",
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

func RenderLoadBalancerMonitorNsxvConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":     "1",
		"Name":               "NSX-V HTTP Monitor",
		"Description":        "An NSX-V HTTP health check monitor",
		"MonitorType":        "http",
		"MonitorInterval":    "10",
		"MonitorTimeout":     "15",
		"MaxRetry":           "3",
		"SendData":           "GET / HTTP/1.0",
		"SendType":           "GET",
		"ReceiveData":        "",
		"ReceiveCode":        "200",
		"MonitorDestination": "/health",
		"MonitorUsername":    "",
		"MonitorPassword":    "",
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
	templatePath := filepath.Join(dir, "example_nsxv.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
