// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package configredact

import (
	"reflect"
	"testing"
)

// TestClusterWithholdsExternalKubernetesCredentials pins MORPH-19991: the keys
// an External Kubernetes cluster's config carries its credentials in never
// reach state, while the rest of the config is untouched.
func TestClusterWithholdsExternalKubernetesCredentials(t *testing.T) {
	t.Parallel()

	in := map[string]any{ //nolint:gosec // G101: illustrative, not a credential
		"apiUrl":            "https://10.227.122.230:6443",
		"apiToken":          "eyJhbGciOiJSUzI1NiIsImtpZCI6Ik…",
		"serviceAccess":     "apiVersion: v1\nkind: Config\n…",
		"cantRecommend":     "true",
		"helmStatusMessage": "Helm is not installed on this cluster",
	}

	out, withheld := Cluster(in)

	want := map[string]any{
		"apiUrl":            "https://10.227.122.230:6443",
		"cantRecommend":     "true",
		"helmStatusMessage": "Helm is not installed on this cluster",
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("redacted config = %#v, want %#v", out, want)
	}

	if wantW := []string{"apiToken", "serviceAccess"}; !reflect.DeepEqual(withheld, wantW) {
		t.Errorf("withheld = %v, want %v", withheld, wantW)
	}

	// The input is not modified.
	if _, ok := in["apiToken"]; !ok {
		t.Error("input map was mutated")
	}
}

// TestClusterKeyRules covers the exact list, the name pattern, and names the
// pattern must leave alone.
func TestClusterKeyRules(t *testing.T) {
	t.Parallel()

	cases := []struct {
		key      string
		withheld bool
	}{
		// Exact keys, whatever the case of their suffix.
		{"apiToken", true},
		{"serviceAccess", true},
		{"serviceToken", true},
		{"kubeToken", true},
		{"kubeConfig", true},
		{"joinCommand", true},
		{"joincommand", true},
		{"APITOKEN", true},
		// The name says it is a credential.
		{"bearerToken", true},
		{"adminPassword", true},
		{"sshPasswd", true},
		{"clientSecret", true},
		{"keyPassphrase", true},
		{"sshPrivateKey", true},
		{"kubeconfig", true},
		{"KUBECONFIG", true},
		// Names that only contain the word are kept: they are settings about
		// a credential, not the credential.
		{"tokenExpiry", false},
		{"passwordPolicy", false},
		{"secretsBackend", false},
		// Ordinary configuration.
		{"apiUrl", false},
		{"cpuArch", false},
		{"dns1", false},
		{"virtSwitchName", false},
		{"podCidr", false},
	}

	for _, c := range cases {
		out, _ := Cluster(map[string]any{c.key: "v"})
		_, kept := out[c.key]

		if kept == c.withheld {
			t.Errorf("%q: withheld = %v, want %v", c.key, !kept, c.withheld)
		}
	}
}

// TestClusterRedactsNested pins that a credential inside a nested map or a
// list of maps is withheld too, as the appliance's own redaction does.
func TestClusterRedactsNested(t *testing.T) {
	t.Parallel()

	in := map[string]any{
		"name": "x",
		"registry": map[string]any{
			"url":      "https://registry",
			"password": "hunter2",
		},
		"nodes": []any{
			map[string]any{"name": "n1", "apiToken": "t1"},
			"plain",
		},
	}

	out, withheld := Cluster(in)

	want := map[string]any{
		"name":     "x",
		"registry": map[string]any{"url": "https://registry"},
		"nodes":    []any{map[string]any{"name": "n1"}, "plain"},
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("redacted config = %#v, want %#v", out, want)
	}

	if wantW := []string{"apiToken", "password"}; !reflect.DeepEqual(withheld, wantW) {
		t.Errorf("withheld = %v, want %v", withheld, wantW)
	}
}

// TestClusterRedactsTypedShapes covers the shapes a caller might build by
// hand rather than decode from JSON.
func TestClusterRedactsTypedShapes(t *testing.T) {
	t.Parallel()

	in := map[string]any{
		"labels": map[string]string{"env": "lab", "apiToken": "t"},
		"hosts":  []map[string]any{{"name": "h1", "sshPassword": "p"}},
	}

	out, withheld := Cluster(in)

	want := map[string]any{
		"labels": map[string]string{"env": "lab"},
		"hosts":  []map[string]any{{"name": "h1"}},
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("redacted config = %#v, want %#v", out, want)
	}

	if wantW := []string{"apiToken", "sshPassword"}; !reflect.DeepEqual(withheld, wantW) {
		t.Errorf("withheld = %v, want %v", withheld, wantW)
	}
}

func TestClusterNil(t *testing.T) {
	t.Parallel()

	out, withheld := Cluster(nil)
	if out != nil || withheld != nil {
		t.Errorf("nil config: got %v, %v", out, withheld)
	}
}
