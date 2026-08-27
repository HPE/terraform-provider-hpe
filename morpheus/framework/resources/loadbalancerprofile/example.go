// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_profile/example_http.tf example_http.tf.tmpl LoadBalancerId "1" Name "HTTP Profile" Description "Example HTTP profile" HttpIdleTimeout "15" RequestHeaderSize "1024" HttpsRedirect "true" XForwardedFor "INSERT"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_profile/example_cookie.tf example_cookie.tf.tmpl LoadBalancerId "1" Name "Cookie Profile" CookieMode "INSERT" CookieType "session" CookieName "SERVERID"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_profile/example_client_ssl.tf example_client_ssl.tf.tmpl LoadBalancerId "1" Name "Client SSL Profile" SslSuite "CUSTOM"
//go:generate ../../../../bin/render -out examples/resources/morpheus_load_balancer_profile/example_minimal.tf example_minimal.tf.tmpl LoadBalancerId "1" Name "Fast TCP Profile"

func RenderLoadBalancerProfileHTTPConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId":    "1",
		"Name":              "HTTP Profile",
		"Description":       "Example HTTP profile",
		"HttpIdleTimeout":   "15",
		"RequestHeaderSize": "1024",
		"HttpsRedirect":     "true",
		"XForwardedFor":     "INSERT",
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
	templatePath := filepath.Join(dir, "example_http.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderLoadBalancerProfileCookieConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId": "1",
		"Name":           "Cookie Profile",
		"CookieMode":     "INSERT",
		"CookieType":     "session",
		"CookieName":     "SERVERID",
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
	templatePath := filepath.Join(dir, "example_cookie.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderLoadBalancerProfileClientSSLConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId": "1",
		"Name":           "Client SSL Profile",
		"SslSuite":       "CUSTOM",
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
	templatePath := filepath.Join(dir, "example_client_ssl.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderLoadBalancerProfileMinimalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId": "1",
		"Name":           "Fast TCP Profile",
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
	templatePath := filepath.Join(dir, "example_minimal.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}
