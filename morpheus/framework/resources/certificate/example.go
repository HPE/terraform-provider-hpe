// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package certificate

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_certificate/example.tf example.tf.tmpl Name "wildcard-example-com" CertFile file("${path.module}/certs/wildcard.crt") KeyFile file("${path.module}/certs/wildcard.key") DomainName "*.example.com" Description "Wildcard certificate for example.com"

func RenderCertificateConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "wildcard-example-com",
		"CertFile":    "\"-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----\"",
		"KeyFile":     "\"-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----\"",
		"DomainName":  "*.example.com",
		"Description": "Wildcard certificate for example.com",
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
