// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// writeJar writes content to a temp .jar file and returns its path and digest.
func writeJar(t *testing.T, content string) (string, string) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "p.jar")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write jar: %v", err)
	}

	sum := sha256.Sum256([]byte(content))

	return path, hex.EncodeToString(sum[:])
}

func TestUnitHashKnownJarPath(t *testing.T) {
	t.Parallel()

	path, want := writeJar(t, "hello-jar")

	tests := []struct {
		name      string
		plan      PluginModel
		wantValue types.String
	}{
		{
			name:      "null path is unknown",
			plan:      PluginModel{JarPath: types.StringNull()},
			wantValue: types.StringUnknown(),
		},
		{
			name:      "unknown path is unknown",
			plan:      PluginModel{JarPath: types.StringUnknown()},
			wantValue: types.StringUnknown(),
		},
		{
			name:      "missing file is unknown",
			plan:      PluginModel{JarPath: types.StringValue("/no/such/plugin.jar")},
			wantValue: types.StringUnknown(),
		},
		{
			name:      "readable file hashes to digest",
			plan:      PluginModel{JarPath: types.StringValue(path)},
			wantValue: types.StringValue(want),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := hashKnownJarPath(tc.plan)
			if !got.Equal(tc.wantValue) {
				t.Fatalf("hashKnownJarPath() = %v, want %v", got, tc.wantValue)
			}
		})
	}
}

func TestUnitHashFile(t *testing.T) {
	t.Parallel()

	path, want := writeJar(t, "contents")

	got, err := hashFile(path)
	if err != nil {
		t.Fatalf("hashFile() error = %v", err)
	}
	if got != want {
		t.Fatalf("hashFile() = %s, want %s", got, want)
	}

	if _, err := hashFile("/no/such/file.jar"); err == nil {
		t.Fatal("hashFile() on missing file: expected error, got nil")
	}
}

// TestUnitJarSha256ConfiguredValueWins verifies that when jar_sha256 is set in
// config the plan modifier echoes it verbatim and never touches state/plan.
func TestUnitJarSha256ConfiguredValueWins(t *testing.T) {
	t.Parallel()

	configured := types.StringValue("deadbeef")

	req := planmodifier.StringRequest{
		ConfigValue: configured,
	}
	resp := &planmodifier.StringResponse{PlanValue: types.StringUnknown()}

	jarSha256PlanModifier{}.PlanModifyString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.PlanValue.Equal(configured) {
		t.Fatalf("PlanValue = %v, want %v", resp.PlanValue, configured)
	}
}
