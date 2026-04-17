// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ./generate_example.sh

func RenderKeyValueConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "keyValue-input",
		"OptionTypeDefaultValue":          `jsonencode([{ key = "a", value = "b" }, { key = "c", value = "d" }])`,
		"OptionTypeDescription":           "Terraform keyValue example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "KeyValue",
		"OptionTypeFieldName":             "keyValue",
		"OptionTypeHelpBlock":             "Select a key-value pair",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf keyValue example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "keyValue",
		"OptionTypeConvertToObject":       "true",
		"OptionTypeKeyPlaceholder":        "Key123",
		"OptionTypeValuePlaceholder":      "Value123",
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
	templatePath := filepath.Join(dir, "form_key_value.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderVirtualImageConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "virtual-image",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform virtual-image example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Virtual Image",
		"OptionTypeFieldName":             "virtual-image",
		"OptionTypeHelpBlock":             "Select a virtual image",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf virtual-image example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "virtual-image",
		"OptionTypeCloudFieldType":        "id",
		"OptionTypeCloudId":               "1",
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
	templatePath := filepath.Join(dir, "form_virtual_image.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderVmwFoldersConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "vmw-folders-input",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform vmwFolders example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "VmwFolders",
		"OptionTypeFieldName":             "vmwFolders",
		"OptionTypeHelpBlock":             "Select a vmwFolder",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf vmwFolders example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "vmwFolders",
		"OptionTypeGroupFieldType":        "value",
		"OptionTypeGroupId":               "1",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
		"OptionTypePlanFieldType":         "value",
		"OptionTypePlanId":                "1",
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
	templatePath := filepath.Join(dir, "form_vmw_folders.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderFileContentConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "fileContent",
		"OptionTypeDescription":           "Terraform fileContent example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "FileContent",
		"OptionTypeFieldName":             "fileContent",
		"OptionTypeHelpBlock":             "Set fileContent",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf fileContent example",
		"OptionTypePlaceholder":           "testing123",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "fileContent",
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
	templatePath := filepath.Join(dir, "form_file_content.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderSelectConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "select-input",
		"OptionTypeDefaultValue":          `"level1"`,
		"OptionTypeDescription":           "Terraform select example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Select Test",
		"OptionTypeFieldName":             "selectTest",
		"OptionTypeHelpBlock":             "Select an option",
		"OptionTypeHidden":                "true",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf example select",
		"OptionTypeOptionListId":          "1",
		"OptionTypePlaceholder":           "Testing 123",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "select",
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
	templatePath := filepath.Join(dir, "form_select.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderRadioConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "radio-input",
		"OptionTypeDefaultValue":          "Demo123",
		"OptionTypeDescription":           "Terraform radio example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Radio Test",
		"OptionTypeFieldName":             "radioTest",
		"OptionTypeHelpBlock":             "Select an option",
		"OptionTypeHidden":                "true",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf radio example",
		"OptionTypeOptionListId":          "1",
		"OptionTypePlaceholder":           "Testing 123",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "radio",
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
	templatePath := filepath.Join(dir, "form_radio.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderTextConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "test-input",
		"OptionTypeDefaultValue":          "Demo123",
		"OptionTypeDescription":           "Terraform text example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Testin",
		"OptionTypeFieldName":             "test",
		"OptionTypeHelpBlock":             "Help block example",
		"OptionTypeHidden":                "true",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf text example",
		"OptionTypePlaceholder":           "Testing 123",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "text",
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
	templatePath := filepath.Join(dir, "form_text.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderCheckboxConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "checkbox-input",
		"OptionTypeDefaultChecked":        "true",
		"OptionTypeDescription":           "Terraform checkbox example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "checkbox input",
		"OptionTypeFieldName":             "checkboxInput",
		"OptionTypeHelpBlock":             "Help block example",
		"OptionTypeHidden":                "true",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf checkbox example",
		"OptionTypePlaceholder":           "Testing 123",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "checkbox",
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
	templatePath := filepath.Join(dir, "form_checkbox.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderHiddenConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "hidden-input",
		"OptionTypeDefaultValue":          "test",
		"OptionTypeDescription":           "Terraform hidden input example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "hidden input",
		"OptionTypeFieldName":             "hiddenInput",
		"OptionTypeHelpBlock":             "Help block example",
		"OptionTypeHidden":                "true",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf hidden input example",
		"OptionTypePlaceholder":           "Testing 123",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "hidden",
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
	templatePath := filepath.Join(dir, "form_hidden.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderHTTPHeaderConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "httpheader-input",
		"OptionTypeDefaultValue":          "[{ name = \"header1\", value = \"value1\", masked = false }]",
		"OptionTypeDescription":           "Terraform HTTP header input example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "HTTP Headers",
		"OptionTypeFieldName":             "httpHeaders",
		"OptionTypeHelpBlock":             "Configure HTTP headers",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf httpheader example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "httpHeader",
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
	templatePath := filepath.Join(dir, "form_httpheader.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderNumberConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "number-input",
		"OptionTypeDefaultValue":          "4",
		"OptionTypeDescription":           "Terraform number example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "number input",
		"OptionTypeFieldName":             "numberInput",
		"OptionTypeHelpBlock":             "Help block example",
		"OptionTypeHidden":                "true",
		"OptionTypeLocked":                "true",
		"OptionTypeMaxValue":              "44",
		"OptionTypeMinValue":              "3",
		"OptionTypeName":                  "tf number input example",
		"OptionTypePlaceholder":           "Testing 123",
		"OptionTypeRequired":              "true",
		"OptionTypeStep":                  "2",
		"OptionTypeType":                  "number",
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
	templatePath := filepath.Join(dir, "form_number.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderNetworkManagerConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                     "demo",
		"Description":              "demo",
		"Labels":                   "[\"terraform\", \"demo\"]",
		"Name":                     "demo",
		"OptionTypeCode":           "network-manager-input",
		"OptionTypeCloudFieldType": "value",
		"OptionTypeCloudId":        "1",
		"OptionTypeDefaultValue": `jsonencode([{ primaryInterface = true, displayOrder = 0, ipMode = "",` +
			`ipAddress = "", networkInterfaceTypeId = "4",` +
			`network = { id = "network-216", pool = "{id: \"\"}" } },` +
			`{ primaryInterface = false, displayOrder = 1, ipMode = "dhcp", ipAddress = "",` +
			`networkInterfaceTypeId = 4, network = { id = "network-216", pool = "{id: \"\"}" } }])`,
		"OptionTypeDescription":              "Terraform network manager example",
		"OptionTypeDisplayValueOnDetails":    "true",
		"OptionTypeEnableIPModeSelection":    "true",
		"OptionTypeExcludeFromSearch":        "true",
		"OptionTypeExportMeta":               "true",
		"OptionTypeFieldLabel":               "network input",
		"OptionTypeFieldName":                "networkInput",
		"OptionTypeGroupFieldType":           "value",
		"OptionTypeGroupId":                  "1",
		"OptionTypeHelpBlock":                "Select a network",
		"OptionTypeHidden":                   "false",
		"OptionTypeLayoutFieldType":          "value",
		"OptionTypeLayoutId":                 "1",
		"OptionTypeLocked":                   "true",
		"OptionTypeName":                     "tf network manager example",
		"OptionTypePlaceholder":              "Select network",
		"OptionTypePoolFieldType":            "value",
		"OptionTypePoolId":                   "1",
		"OptionTypeRequired":                 "true",
		"OptionTypeShowNetworkTypeSelection": "true",
		"OptionTypeType":                     "networkManager",
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
	templatePath := filepath.Join(dir, "form_network_manager.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderCloudConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "cloud-input",
		"OptionTypeCloudType":             "4",
		"OptionTypeDefaultValue":          "test123",
		"OptionTypeDescription":           "Terraform cloud example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "cloud input",
		"OptionTypeFieldName":             "cloudInput",
		"OptionTypeFilterFromResource":    "true",
		"OptionTypeGroupFieldType":        "value",
		"OptionTypeGroupId":               "1",
		"OptionTypeHelpBlock":             "Select a cloud",
		"OptionTypeHidden":                "false",
		"OptionTypeInstanceTypeCode":      "apache",
		"OptionTypeInstanceTypeFieldType": "value",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf cloud example",
		"OptionTypePlaceholder":           "Select cloud",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "cloud",
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
	templatePath := filepath.Join(dir, "form_cloud.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderLayoutConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "layout-input",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform layout example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "layout input",
		"OptionTypeFieldName":             "layoutInput",
		"OptionTypeGroupFieldType":        "value",
		"OptionTypeGroupId":               "1",
		"OptionTypeHelpBlock":             "Select a layout",
		"OptionTypeHidden":                "false",
		"OptionTypeInstanceTypeCode":      "apache",
		"OptionTypeInstanceTypeFieldType": "value",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf layout example",
		"OptionTypePlaceholder":           "Select layout",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "layout",
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
	templatePath := filepath.Join(dir, "form_layout.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeAllowReadOnly":         "true",
		"OptionTypeCode":                  "group-input",
		"OptionTypeDefaultValue":          "test123",
		"OptionTypeDescription":           "Terraform group example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "group input",
		"OptionTypeFieldName":             "groupInput",
		"OptionTypeHelpBlock":             "Select a group",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf group example",
		"OptionTypePlaceholder":           "Select group",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "group",
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
	templatePath := filepath.Join(dir, "form_group.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderDiskManagerConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                     "demo",
		"Description":              "demo",
		"Labels":                   "[\"terraform\", \"demo\"]",
		"Name":                     "demo",
		"OptionTypeCloudFieldType": "value",
		"OptionTypeCloudId":        "1",
		"OptionTypeCode":           "disk-manager-input",
		"OptionTypeDefaultValue": `jsonencode([{ rootVolume = true, name = "root", size = 10,` +
			`sizeBytes = 10737418240, minStorage = 0, displayOrder = 0, storageType = 1, datastoreId = "52" },` +
			`{ rootVolume = false, name = "data-1", size = 20, sizeBytes = 21474836480,` +
			`minStorage = 0, displayOrder = 1, datastoreId = "autoCluster", storageType = 1 }])`,
		"OptionTypeDescription":                "Terraform disk manager example",
		"OptionTypeDisplayValueOnDetails":      "true",
		"OptionTypeEnableDatastoreSelection":   "true",
		"OptionTypeEnableDiskTypeSelection":    "true",
		"OptionTypeEnableStorageTypeSelection": "true",
		"OptionTypeExcludeFromSearch":          "true",
		"OptionTypeExportMeta":                 "true",
		"OptionTypeFieldLabel":                 "disk manager input",
		"OptionTypeFieldName":                  "diskManagerInput",
		"OptionTypeGroupFieldType":             "value",
		"OptionTypeGroupId":                    "1",
		"OptionTypeHelpBlock":                  "Configure disks",
		"OptionTypeHidden":                     "false",
		"OptionTypeLayoutFieldType":            "value",
		"OptionTypeLayoutId":                   "1",
		"OptionTypeLocked":                     "true",
		"OptionTypeName":                       "tf disk manager example",
		"OptionTypePlanFieldType":              "value",
		"OptionTypePlanId":                     "1",
		"OptionTypePoolFieldType":              "value",
		"OptionTypePoolId":                     "1",
		"OptionTypeRequired":                   "true",
		"OptionTypeType":                       "diskManager",
		"OptionTypeVirtualImageFieldType":      "value",
		"OptionTypeImageId":                    "1",
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
	templatePath := filepath.Join(dir, "form_disk_manager.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderPlanConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                     "demo",
		"Description":              "demo",
		"Labels":                   "[\"terraform\", \"demo\"]",
		"Name":                     "demo",
		"OptionTypeCloudFieldType": "value",
		"OptionTypeCloudId":        "1",
		"OptionTypeCode":           "plan-input",
		"OptionTypeDefaultValue": `jsonencode({ id = 1088, maxMemory = 8589934592,` +
			`maxCores = "4", coresPerSocket = "2" })`,
		"OptionTypeDescription":           "Terraform plan example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "plan input",
		"OptionTypeFieldName":             "planInput",
		"OptionTypeGroupFieldType":        "value",
		"OptionTypeGroupId":               "1",
		"OptionTypeHelpBlock":             "Select a plan",
		"OptionTypeHidden":                "false",
		"OptionTypeLayoutFieldType":       "value",
		"OptionTypeLayoutId":              "1",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf plan example",
		"OptionTypePlaceholder":           "Select plan",
		"OptionTypePoolFieldType":         "value",
		"OptionTypePoolId":                "1",
		"OptionTypeRequired":              "true",
		"OptionTypeShowPricing":           "false",
		"OptionTypeType":                  "plan",
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
	templatePath := filepath.Join(dir, "form_plan.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderFieldGroupsConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                                       "demo",
		"Description":                                "demo",
		"FieldGroup1CollapsedByDefault":              "true",
		"FieldGroup1Collapsible":                     "true",
		"FieldGroup1Description":                     "testin",
		"FieldGroup1Name":                            "fg1",
		"FieldGroup1OptionTypeCode":                  "test-input",
		"FieldGroup1OptionTypeDefaultValue":          "Demo123",
		"FieldGroup1OptionTypeDescription":           "Terraform text input example",
		"FieldGroup1OptionTypeDisplayValueOnDetails": "true",
		"FieldGroup1OptionTypeExcludeFromSearch":     "true",
		"FieldGroup1OptionTypeExportMeta":            "true",
		"FieldGroup1OptionTypeFieldLabel":            "Testing 1",
		"FieldGroup1OptionTypeFieldName":             "test1",
		"FieldGroup1OptionTypeHelpBlock":             "Help block example",
		"FieldGroup1OptionTypeHidden":                "false",
		"FieldGroup1OptionTypeLocked":                "true",
		"FieldGroup1OptionTypeName":                  "tf field group 1 text input example",
		"FieldGroup1OptionTypePlaceholder":           "Testing 123",
		"FieldGroup1OptionTypeRequired":              "true",
		"FieldGroup1OptionTypeType":                  "text",
		"FieldGroup2CollapsedByDefault":              "true",
		"FieldGroup2Collapsible":                     "true",
		"FieldGroup2Description":                     "testin",
		"FieldGroup2Name":                            "fg2",
		"FieldGroup2OptionTypeCode":                  "test-input",
		"FieldGroup2OptionTypeDefaultValue":          "Demo123",
		"FieldGroup2OptionTypeDescription":           "Terraform text input example",
		"FieldGroup2OptionTypeDisplayValueOnDetails": "true",
		"FieldGroup2OptionTypeExcludeFromSearch":     "true",
		"FieldGroup2OptionTypeExportMeta":            "true",
		"FieldGroup2OptionTypeFieldLabel":            "Testing 2",
		"FieldGroup2OptionTypeFieldName":             "test2",
		"FieldGroup2OptionTypeHelpBlock":             "Help block example",
		"FieldGroup2OptionTypeHidden":                "false",
		"FieldGroup2OptionTypeLocked":                "true",
		"FieldGroup2OptionTypeName":                  "tf field group 2 text input example",
		"FieldGroup2OptionTypePlaceholder":           "Testing 123",
		"FieldGroup2OptionTypeRequired":              "true",
		"FieldGroup2OptionTypeType":                  "text",
		"Labels":                                     "[\"terraform\", \"demo\"]",
		"Name":                                       "demo",
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
	templatePath := filepath.Join(dir, "form_field_groups.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderEnvironmentConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "environment-input",
		"OptionTypeDefaultValue":          "staging",
		"OptionTypeDescription":           "Terraform environment example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Environment",
		"OptionTypeFieldName":             "environment",
		"OptionTypeHelpBlock":             "Select an environment",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf environment example",
		"OptionTypePlaceholder":           "",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "environment",
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
	templatePath := filepath.Join(dir, "form_environment.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderServersInputConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "servers-input",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform servers-input example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Server",
		"OptionTypeFieldName":             "server",
		"OptionTypeHelpBlock":             "Select a server",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf servers-input example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "servers-input",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
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
	templatePath := filepath.Join(dir, "form_servers_input.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderResourcePoolConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "resource-pool-input",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform resourcePool example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Resource Pool",
		"OptionTypeFieldName":             "resourcePool",
		"OptionTypeHelpBlock":             "Select a resource pool",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf resourcePool example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "resourcePool",
		"OptionTypeGroupFieldType":        "value",
		"OptionTypeGroupId":               "1",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
		"OptionTypePlanFieldType":         "value",
		"OptionTypePlanId":                "1",
		"OptionTypeLayoutFieldType":       "value",
		"OptionTypeLayoutId":              "1",
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
	templatePath := filepath.Join(dir, "form_resource_pool.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderSecGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "sec-group-input",
		"OptionTypeDefaultValue":          `jsonencode([{ id = "sec-group-default" }])`,
		"OptionTypeDescription":           "Terraform secGroup example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Security Groups",
		"OptionTypeFieldName":             "securityGroups",
		"OptionTypeHelpBlock":             "Select security groups",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf secGroup example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "secGroup",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
		"OptionTypePoolField":             "resourcePool",
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
	templatePath := filepath.Join(dir, "form_sec_group.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderTagConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "tag-input",
		"OptionTypeDefaultValue":          `jsonencode([{ name = "Sample Name", value = "Sample Value" }])`,
		"OptionTypeDescription":           "Terraform tag example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Tags",
		"OptionTypeFieldName":             "tags",
		"OptionTypeHelpBlock":             "Configure tags",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf tag example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "tag",
		"OptionTypeGroupFieldType":        "value",
		"OptionTypeGroupId":               "1",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
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
	templatePath := filepath.Join(dir, "form_tag.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderInstancesInputConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "instances-input",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform instances-input example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Instance",
		"OptionTypeFieldName":             "instance",
		"OptionTypeHelpBlock":             "Select an instance",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf instances-input example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "instances-input",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
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
	templatePath := filepath.Join(dir, "form_servers_input.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderPortsConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":           "demo",
		"Description":    "demo",
		"Labels":         "[\"terraform\", \"demo\"]",
		"Name":           "demo",
		"OptionTypeCode": "ports-input",
		"OptionTypeDefaultValue": `jsonencode([{ name = "standard", externalPort = "80", loadBalanceProtocol = "HTTP" },` +
			`{ name = "ssl-title", externalPort = "443", loadBalanceProtocol = "HTTPS" },` +
			`{ name = "tcp", externalPort = "40", loadBalanceProtocol = "TCP" }])`,
		"OptionTypeDescription":           "Terraform ports example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Exposed Ports",
		"OptionTypeFieldName":             "ports",
		"OptionTypeHelpBlock":             "Configure exposed ports",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf ports example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "ports",
		"OptionTypeGroupField":            "myGroup",
		"OptionTypeCloudField":            "myCloud",
		"OptionTypeLayoutField":           "myLayout",
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
	templatePath := filepath.Join(dir, "form_ports.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderLogoSelectorConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":           "demo",
		"Description":    "demo",
		"Labels":         "[\"terraform\", \"demo\"]",
		"Name":           "demo",
		"OptionTypeCode": "logo-selector-input",
		"OptionTypeDefaultValue": `jsonencode({ value = "identicon",` +
			`settings = { type = "identicon", iconLabel = "example" } })`,
		"OptionTypeDescription":           "Terraform logo selector example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Select Logo",
		"OptionTypeFieldName":             "logoSelector",
		"OptionTypeHelpBlock":             "Select or upload a logo",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf logo selector example",
		"OptionTypePlaceholder":           "",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "logoSelector",
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
	templatePath := filepath.Join(dir, "form_logo_selector.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderByteSizeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "bytesize-input",
		"OptionTypeDefaultValue":          "GB",
		"OptionTypeDescription":           "Terraform byteSize example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeDisplay":               "48318382080",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Byte Size",
		"OptionTypeFieldName":             "byteSize",
		"OptionTypeHelpBlock":             "Select byte size display",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeLockDisplay":           "false",
		"OptionTypeName":                  "tf byteSize example",
		"OptionTypePlaceholder":           "",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "byteSize",
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
	templatePath := filepath.Join(dir, "form_bytesize.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderCodeEditorConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "code-editor-input",
		"OptionTypeCodeLanguage":          "bash",
		"OptionTypeDefaultValue":          "echo \"hello world\"",
		"OptionTypeDescription":           "Terraform code-editor example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Code Editor",
		"OptionTypeFieldName":             "codeEditor",
		"OptionTypeHelpBlock":             "Enter code",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf code-editor example",
		"OptionTypePlaceholder":           "",
		"OptionTypeRequired":              "true",
		"OptionTypeShowLineNumbers":       "true",
		"OptionTypeType":                  "code-editor",
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
	templatePath := filepath.Join(dir, "form_code_editor.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderPasswordConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeAllowPasswordPeek":     "true",
		"OptionTypeCode":                  "password-input",
		"OptionTypeDefaultValue":          "",
		"OptionTypeDescription":           "Terraform password example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Password",
		"OptionTypeFieldName":             "password",
		"OptionTypeHelpBlock":             "Enter a secure password",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf password example",
		"OptionTypePlaceholder":           "Enter password",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "password",
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
	templatePath := filepath.Join(dir, "form_password.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderTextAreaConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "textarea-input",
		"OptionTypeDefaultValue":          "Sample text",
		"OptionTypeDescription":           "Terraform textarea example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Text Area",
		"OptionTypeFieldName":             "textArea",
		"OptionTypeHelpBlock":             "Enter multiple lines of text",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf textarea example",
		"OptionTypePlaceholder":           "Enter text",
		"OptionTypeRequired":              "true",
		"OptionTypeTextRows":              "5",
		"OptionTypeType":                  "textarea",
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
	templatePath := filepath.Join(dir, "form_textarea.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderTextArrayConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "text-array-input",
		"OptionTypeDefaultValue":          "jsonencode([\"item1\", \"item2\", \"item3\"])",
		"OptionTypeDelimiter":             ",",
		"OptionTypeDescription":           "Terraform textArray example",
		"OptionTypeDisplayValueOnDetails": "true",
		"OptionTypeExcludeFromSearch":     "true",
		"OptionTypeExportMeta":            "true",
		"OptionTypeFieldLabel":            "Text Array",
		"OptionTypeFieldName":             "textArray",
		"OptionTypeHelpBlock":             "Enter comma-separated values",
		"OptionTypeHidden":                "false",
		"OptionTypeLocked":                "true",
		"OptionTypeName":                  "tf textArray example",
		"OptionTypeRequired":              "true",
		"OptionTypeType":                  "textArray",
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
	templatePath := filepath.Join(dir, "form_text_array.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderTypeaheadConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                              "demo",
		"Description":                       "demo",
		"Labels":                            "[\"terraform\", \"demo\"]",
		"Name":                              "demo",
		"OptionTypeAllowDuplicates":         "false",
		"OptionTypeAllowMultipleSelections": "false",
		"OptionTypeCode":                    "typeahead-input",
		"OptionTypeCustomData":              "{}",
		"OptionTypeDefaultValue":            "test",
		"OptionTypeDescription":             "Terraform typeahead example",
		"OptionTypeDisplayValueOnDetails":   "true",
		"OptionTypeExcludeFromSearch":       "true",
		"OptionTypeExportMeta":              "true",
		"OptionTypeFieldLabel":              "Typeahead",
		"OptionTypeFieldName":               "typeahead",
		"OptionTypeHelpBlock":               "Select an option from the list",
		"OptionTypeHidden":                  "false",
		"OptionTypeLocked":                  "true",
		"OptionTypeName":                    "tf typeahead example",
		"OptionTypeOptionListId":            "1",
		"OptionTypePlaceholder":             "Search...",
		"OptionTypeRequired":                "true",
		"OptionTypeSortable":                "true",
		"OptionTypeType":                    "typeahead",
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
	templatePath := filepath.Join(dir, "form_typeahead.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}
