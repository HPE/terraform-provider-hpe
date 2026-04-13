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

// RenderFormConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderFormConfig(t *testing.T, overrides map[string]string) (string, error) {
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
		"FieldGroup1OptionTypeFieldLabel":            "Testin",
		"FieldGroup1OptionTypeFieldName":             "test",
		"FieldGroup1OptionTypeHelpBlock":             "Is this working now",
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
		"FieldGroup2OptionTypeFieldLabel":            "Testin",
		"FieldGroup2OptionTypeFieldName":             "test",
		"FieldGroup2OptionTypeHelpBlock":             "Is this working now",
		"FieldGroup2OptionTypeHidden":                "false",
		"FieldGroup2OptionTypeLocked":                "true",
		"FieldGroup2OptionTypeName":                  "tf field group 2 text input example",
		"FieldGroup2OptionTypePlaceholder":           "Testing 123",
		"FieldGroup2OptionTypeRequired":              "true",
		"FieldGroup2OptionTypeType":                  "text",
		"Labels":                                     "[\"terraform\", \"demo\"]",
		"Name":                                       "demo",
		"OptionType1Code":                            "select-input",
		"OptionType1DefaultValue":                    "test123",
		"OptionType1Description":                     "Terraform select example",
		"OptionType1DisplayValueOnDetails":           "true",
		"OptionType1ExcludeFromSearch":               "true",
		"OptionType1ExportMeta":                      "true",
		"OptionType1FieldLabel":                      "Select Test",
		"OptionType1FieldName":                       "selectTest",
		"OptionType1HelpBlock":                       "Select an option",
		"OptionType1Hidden":                          "true",
		"OptionType1Locked":                          "true",
		"OptionType1Name":                            "tf example select",
		"OptionType1OptionListId":                    "1",
		"OptionType1Placeholder":                     "Testing 123",
		"OptionType1Required":                        "true",
		"OptionType1Type":                            "select",
		"OptionType2Code":                            "radio-input",
		"OptionType2DefaultValue":                    "Demo123",
		"OptionType2Description":                     "Terraform radio example",
		"OptionType2DisplayValueOnDetails":           "true",
		"OptionType2ExcludeFromSearch":               "true",
		"OptionType2ExportMeta":                      "true",
		"OptionType2FieldLabel":                      "Radio Test",
		"OptionType2FieldName":                       "radioTest",
		"OptionType2HelpBlock":                       "Select an option",
		"OptionType2Hidden":                          "true",
		"OptionType2Locked":                          "true",
		"OptionType2Name":                            "tf radio example",
		"OptionType2OptionListId":                    "1",
		"OptionType2Placeholder":                     "Testing 123",
		"OptionType2Required":                        "true",
		"OptionType2Type":                            "radio",
		"OptionType3Code":                            "test-input",
		"OptionType3DefaultValue":                    "Demo123",
		"OptionType3Description":                     "Terraform text example",
		"OptionType3DisplayValueOnDetails":           "true",
		"OptionType3ExcludeFromSearch":               "true",
		"OptionType3ExportMeta":                      "true",
		"OptionType3FieldLabel":                      "Testin",
		"OptionType3FieldName":                       "test",
		"OptionType3HelpBlock":                       "Is this working now",
		"OptionType3Hidden":                          "true",
		"OptionType3Locked":                          "true",
		"OptionType3Name":                            "tf text example",
		"OptionType3Placeholder":                     "Testing 123",
		"OptionType3Required":                        "true",
		"OptionType3Type":                            "text",
		"OptionType4Code":                            "checkbox-input",
		"OptionType4DefaultChecked":                  "true",
		"OptionType4Description":                     "Terraform checkbox example",
		"OptionType4DisplayValueOnDetails":           "true",
		"OptionType4ExcludeFromSearch":               "true",
		"OptionType4ExportMeta":                      "true",
		"OptionType4FieldLabel":                      "checkbox input",
		"OptionType4FieldName":                       "checkboxInput",
		"OptionType4HelpBlock":                       "Is this working now",
		"OptionType4Hidden":                          "true",
		"OptionType4Locked":                          "true",
		"OptionType4Name":                            "tf checkbox example",
		"OptionType4Placeholder":                     "Testing 123",
		"OptionType4Required":                        "true",
		"OptionType4Type":                            "checkbox",
		"OptionType5Code":                            "hidden-input",
		"OptionType5DefaultValue":                    "test",
		"OptionType5Description":                     "Terraform hidden input example",
		"OptionType5DisplayValueOnDetails":           "true",
		"OptionType5ExcludeFromSearch":               "true",
		"OptionType5ExportMeta":                      "true",
		"OptionType5FieldLabel":                      "hidden input",
		"OptionType5FieldName":                       "hiddenInput",
		"OptionType5HelpBlock":                       "Is this working now",
		"OptionType5Hidden":                          "true",
		"OptionType5Locked":                          "true",
		"OptionType5Name":                            "tf hidden input example",
		"OptionType5Placeholder":                     "Testing 123",
		"OptionType5Required":                        "true",
		"OptionType5Type":                            "hidden",
		"OptionType6Code":                            "number-input",
		"OptionType6DefaultValue":                    "4",
		"OptionType6Description":                     "Terraform number example",
		"OptionType6DisplayValueOnDetails":           "true",
		"OptionType6ExcludeFromSearch":               "true",
		"OptionType6ExportMeta":                      "true",
		"OptionType6FieldLabel":                      "number input",
		"OptionType6FieldName":                       "numberInput",
		"OptionType6HelpBlock":                       "Is this working now",
		"OptionType6Hidden":                          "true",
		"OptionType6Locked":                          "true",
		"OptionType6MaxValue":                        "44",
		"OptionType6MinValue":                        "3",
		"OptionType6Name":                            "tf number input example",
		"OptionType6Placeholder":                     "Testing 123",
		"OptionType6Required":                        "true",
		"OptionType6Step":                            "2",
		"OptionType6Type":                            "number",
		"OptionType7Code":                            "network-manager-input",
		"OptionType7CloudFieldType":                  "value",
		"OptionType7CloudId":                         "1",
		"OptionType7DefaultValue":                    "test123",
		"OptionType7Description":                     "Terraform network manager example",
		"OptionType7DisplayValueOnDetails":           "true",
		"OptionType7EnableIPModeSelection":           "true",
		"OptionType7ExcludeFromSearch":               "true",
		"OptionType7ExportMeta":                      "true",
		"OptionType7FieldLabel":                      "network input",
		"OptionType7FieldName":                       "networkInput",
		"OptionType7GroupFieldType":                  "value",
		"OptionType7GroupId":                         "1",
		"OptionType7HelpBlock":                       "Select a network",
		"OptionType7Hidden":                          "false",
		"OptionType7LayoutFieldType":                 "value",
		"OptionType7LayoutId":                        "1",
		"OptionType7Locked":                          "true",
		"OptionType7Name":                            "tf network manager example",
		"OptionType7Placeholder":                     "Select network",
		"OptionType7PoolFieldType":                   "value",
		"OptionType7PoolId":                          "1",
		"OptionType7Required":                        "true",
		"OptionType7ShowNetworkTypeSelection":        "true",
		"OptionType7Type":                            "networkManager",
		"OptionType8Code":                            "cloud-input",
		"OptionType8CloudType":                       "4",
		"OptionType8DefaultValue":                    "test123",
		"OptionType8Description":                     "Terraform cloud example",
		"OptionType8DisplayValueOnDetails":           "true",
		"OptionType8ExcludeFromSearch":               "true",
		"OptionType8ExportMeta":                      "true",
		"OptionType8FieldLabel":                      "cloud input",
		"OptionType8FieldName":                       "cloudInput",
		"OptionType8FilterFromResource":              "true",
		"OptionType8GroupFieldType":                  "value",
		"OptionType8GroupId":                         "1",
		"OptionType8HelpBlock":                       "Select a cloud",
		"OptionType8Hidden":                          "false",
		"OptionType8InstanceTypeCode":                "apache",
		"OptionType8InstanceTypeFieldType":           "value",
		"OptionType8Locked":                          "true",
		"OptionType8Name":                            "tf cloud example",
		"OptionType8Placeholder":                     "Select cloud",
		"OptionType8Required":                        "true",
		"OptionType8Type":                            "cloud",
		"OptionType9Code":                            "layout-input",
		"OptionType9CloudFieldType":                  "value",
		"OptionType9CloudId":                         "1",
		"OptionType9DefaultValue":                    "",
		"OptionType9Description":                     "Terraform layout example",
		"OptionType9DisplayValueOnDetails":           "true",
		"OptionType9ExcludeFromSearch":               "true",
		"OptionType9ExportMeta":                      "true",
		"OptionType9FieldLabel":                      "layout input",
		"OptionType9FieldName":                       "layoutInput",
		"OptionType9GroupFieldType":                  "value",
		"OptionType9GroupId":                         "1",
		"OptionType9HelpBlock":                       "Select a layout",
		"OptionType9Hidden":                          "false",
		"OptionType9InstanceTypeCode":                "apache",
		"OptionType9InstanceTypeFieldType":           "value",
		"OptionType9Locked":                          "true",
		"OptionType9Name":                            "tf layout example",
		"OptionType9Placeholder":                     "Select layout",
		"OptionType9Required":                        "true",
		"OptionType9Type":                            "layout",
		"OptionType10AllowReadOnly":                  "true",
		"OptionType10Code":                           "group-input",
		"OptionType10DefaultValue":                   "test123",
		"OptionType10Description":                    "Terraform group example",
		"OptionType10DisplayValueOnDetails":          "true",
		"OptionType10ExcludeFromSearch":              "true",
		"OptionType10ExportMeta":                     "true",
		"OptionType10FieldLabel":                     "group input",
		"OptionType10FieldName":                      "groupInput",
		"OptionType10HelpBlock":                      "Select a group",
		"OptionType10Hidden":                         "false",
		"OptionType10Locked":                         "true",
		"OptionType10Name":                           "tf group example",
		"OptionType10Placeholder":                    "Select group",
		"OptionType10Required":                       "true",
		"OptionType10Type":                           "group",
		"OptionType11CloudFieldType":                 "value",
		"OptionType11CloudId":                        "1",
		"OptionType11Code":                           "disk-manager-input",
		"OptionType11Description":                    "Terraform disk manager example",
		"OptionType11DisplayValueOnDetails":          "true",
		"OptionType11EnableDatastoreSelection":       "true",
		"OptionType11EnableDiskTypeSelection":        "true",
		"OptionType11EnableStorageTypeSelection":     "true",
		"OptionType11ExcludeFromSearch":              "true",
		"OptionType11ExportMeta":                     "true",
		"OptionType11FieldLabel":                     "disk manager input",
		"OptionType11FieldName":                      "diskManagerInput",
		"OptionType11GroupFieldType":                 "value",
		"OptionType11GroupId":                        "1",
		"OptionType11HelpBlock":                      "Configure disks",
		"OptionType11Hidden":                         "false",
		"OptionType11LayoutFieldType":                "value",
		"OptionType11LayoutId":                       "1",
		"OptionType11Locked":                         "true",
		"OptionType11Name":                           "tf disk manager example",
		"OptionType11PlanFieldType":                  "value",
		"OptionType11PlanId":                         "1",
		"OptionType11PoolFieldType":                  "value",
		"OptionType11PoolId":                         "1",
		"OptionType11Required":                       "true",
		"OptionType11Type":                           "diskManager",
		"OptionType11VirtualImageFieldType":          "value",
		"OptionType11ImageId":                        "1",
		"OptionType12CloudFieldType":                 "value",
		"OptionType12CloudId":                        "1",
		"OptionType12Code":                           "plan-input",
		"OptionType12DefaultValue":                   "",
		"OptionType12Description":                    "Terraform plan example",
		"OptionType12DisplayValueOnDetails":          "true",
		"OptionType12ExcludeFromSearch":              "true",
		"OptionType12ExportMeta":                     "true",
		"OptionType12FieldLabel":                     "plan input",
		"OptionType12FieldName":                      "planInput",
		"OptionType12GroupFieldType":                 "value",
		"OptionType12GroupId":                        "1",
		"OptionType12HelpBlock":                      "Select a plan",
		"OptionType12Hidden":                         "false",
		"OptionType12LayoutFieldType":                "value",
		"OptionType12LayoutId":                       "1",
		"OptionType12Locked":                         "true",
		"OptionType12Name":                           "tf plan example",
		"OptionType12Placeholder":                    "Select plan",
		"OptionType12PoolFieldType":                  "value",
		"OptionType12PoolId":                         "1",
		"OptionType12Required":                       "true",
		"OptionType12ShowPricing":                    "false",
		"OptionType12Type":                           "plan",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "form_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderKeyValueConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "keyValue-input",
		"OptionTypeDefaultValue":          "",
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
		"OptionTypeDefaultValue":          "test123",
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
		"OptionTypeHelpBlock":             "Is this working now",
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
		"OptionTypeHelpBlock":             "Is this working now",
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
		"OptionTypeHelpBlock":             "Is this working now",
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
		"OptionTypeHelpBlock":             "Is this working now",
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
		"Code":                               "demo",
		"Description":                        "demo",
		"Labels":                             "[\"terraform\", \"demo\"]",
		"Name":                               "demo",
		"OptionTypeCode":                     "network-manager-input",
		"OptionTypeCloudFieldType":           "value",
		"OptionTypeCloudId":                  "1",
		"OptionTypeDefaultValue":             "test123",
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
		"Code":                                 "demo",
		"Description":                          "demo",
		"Labels":                               "[\"terraform\", \"demo\"]",
		"Name":                                 "demo",
		"OptionTypeCloudFieldType":             "value",
		"OptionTypeCloudId":                    "1",
		"OptionTypeCode":                       "disk-manager-input",
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
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCloudFieldType":        "value",
		"OptionTypeCloudId":               "1",
		"OptionTypeCode":                  "plan-input",
		"OptionTypeDefaultValue":          "",
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
		"FieldGroup1OptionTypeFieldLabel":            "Testin",
		"FieldGroup1OptionTypeFieldName":             "test",
		"FieldGroup1OptionTypeHelpBlock":             "Is this working now",
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
		"FieldGroup2OptionTypeFieldLabel":            "Testin",
		"FieldGroup2OptionTypeFieldName":             "test",
		"FieldGroup2OptionTypeHelpBlock":             "Is this working now",
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
		"OptionTypeDefaultValue":          "",
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
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "ports-input",
		"OptionTypeDefaultValue":          "",
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
		"Code":                            "demo",
		"Description":                     "demo",
		"Labels":                          "[\"terraform\", \"demo\"]",
		"Name":                            "demo",
		"OptionTypeCode":                  "logo-selector-input",
		"OptionTypeDefaultValue":          "identicon",
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
