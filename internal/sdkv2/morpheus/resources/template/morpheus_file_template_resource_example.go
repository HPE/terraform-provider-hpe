// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_file_template/resource.tf morpheus_file_template_resource.tf.tmpl Name tf-terraform-file-template Labels ["demo","template","terraform"] FileName tfcustom.cnf FilePath /etc/my.cnf.d Phase preProvision FileOwner root SettingName myCnf SettingCategory master

func RenderFileTemplateConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            name,
		"Labels":          `["demo", "template", "terraform"]`,
		"FileName":        "tfcustom.cnf",
		"FilePath":        "/etc/my.cnf.d",
		"Phase":           "preProvision",
		"FileContent":     `"# Test MySQL Configuration\n[mysqld]\ninnodb_buffer_pool_size = 128M"`,
		"FileOwner":       "root",
		"SettingName":     "myCnf",
		"SettingCategory": "master",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_file_template_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Labels", defaults["Labels"],
		"FileName", defaults["FileName"],
		"FilePath", defaults["FilePath"],
		"Phase", defaults["Phase"],
		"FileContent", defaults["FileContent"],
		"FileOwner", defaults["FileOwner"],
		"SettingName", defaults["SettingName"],
		"SettingCategory", defaults["SettingCategory"],
	)
}
