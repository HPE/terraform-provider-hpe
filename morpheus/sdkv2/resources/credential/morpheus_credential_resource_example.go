// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package credential

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_username_password.tf credential_resource_username_password.tf.tmpl Name "tf_example" Description "terraform example" Enabled "true" Type "username-password" Username "admin" Password "password12333"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_api_key.tf credential_resource_api_key.tf.tmpl Name "tf_example_credential_api_key" Description "terraform credential example for api key" Enabled "true" Type "api-key" ApiKey "FIEFMIQNQ"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_access_key_secret.tf credential_resource_access_key_secret.tf.tmpl Name "tf_example_credential_access_key_secret" Description "terraform credential example for access key and secret key" Enabled "true" Type "access-key-secret" AccessKey "FIEFMIQNQ" SecretKey "MFMWEIIEIFENF"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_username_keypair.tf credential_resource_username_keypair.tf.tmpl Name "tf_example_credential_username_keypair" Description "terraform credential example for username key pair" Enabled "true" Type "username-keypair" Username "admin" KeyPairId "22"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_username_api_key.tf credential_resource_username_api_key.tf.tmpl Name "tf_example_credential_username_api_key" Description "terraform credential example for username api key" Enabled "true" Type "username-api-key" Username "admin" ApiKey "MFIEIWEIFINEF"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_tenant_username_keypair.tf credential_resource_tenant_username_keypair.tf.tmpl Name "tf_example_credential_tenant_username_keypair" Description "terraform credential example for tenant username keypair" Enabled "true" Type "tenant-username-keypair" Tenant "tenant123" Username "admin" KeyPairId "22"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_username_password_keypair.tf credential_resource_username_password_keypair.tf.tmpl Name "tf_example_credential_username_password_keypair" Description "terraform credential example for username password key pair" Enabled "true" Type "username-password-keypair" Username "admin" Password "password123" KeyPairId "22"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_email_private_key.tf credential_resource_email_private_key.tf.tmpl Name "tf_example_credential_email_private_key" Description "terraform credential example for email private key" Enabled "true" Type "email-private-key" Email "test@example.local" KeyPairId "33"
//go:generate../../../../bin/render -out examples/resources/morpheus_credential/resource_client_id_secret.tf credential_resource_client_id_secret.tf.tmpl Name "tf_example_credential_client_id_secret" Description "terraform credential example for client id secret" Enabled "true" Type "client-id-secret" ClientId "FIEFMIQNQ" ClientSecret "MMEWMIFINWEINFINE"

func RenderCredentialTenantUsernameKeypairConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for tenant username keypair",
		"Enabled":     "true",
		"Type":        "tenant-username-keypair",
		"Tenant":      "tenant123",
		"Username":    "admin",
		"KeyPairId":   "2",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_tenant_username_keypair.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialClientIdSecretConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description":  "terraform credential example for client id secret",
		"Enabled":      "true",
		"Type":         "client-id-secret",
		"ClientId":     "FIEFMIQNQ",
		"ClientSecret": "MMEWMIFINWEINFINE",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_client_id_secret.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialUsernameApiKeyConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for username api key",
		"Enabled":     "true",
		"Type":        "username-api-key",
		"Username":    "admin",
		"ApiKey":      "MFIEIWEIFINEF",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_username_api_key.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialUsernamePasswordKeypairConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for username password key pair",
		"Enabled":     "true",
		"Type":        "username-password-keypair",
		"Username":    "admin",
		"Password":    "password123",
		"KeyPairId":   "2",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_username_password_keypair.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialAccessKeySecretConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for access key and secret key",
		"Enabled":     "true",
		"Type":        "access-key-secret",
		"AccessKey":   "FIEFMIQNQ",
		"SecretKey":   "MFMWEIIEIFENF",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_access_key_secret.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialApiKeyConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for api key",
		"Enabled":     "true",
		"Type":        "api-key",
		"ApiKey":      "FIEFMIQNQ",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_api_key.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialUsernameKeypairConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for username key pair",
		"Enabled":     "true",
		"Type":        "username-keypair",
		"Username":    "admin",
		"KeyPairId":   "2",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_username_keypair.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialUsernamePasswordConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform example",
		"Enabled":     "true",
		"Type":        "username-password",
		"Username":    "admin",
		"Password":    "password12333",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_username_password.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCredentialEmailPrivateKeyConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for email private key",
		"Enabled":     "true",
		"Type":        "email-private-key",
		"Email":       "test@example.local",
		"KeyPairId":   "2",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "credential_resource_email_private_key.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
