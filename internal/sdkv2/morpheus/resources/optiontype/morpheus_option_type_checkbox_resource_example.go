// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_checkbox/resource.tf morpheus_option_type_checkbox_resource.tf.tmpl Name tfcheckboxexample Description "Terraform checkbox option type example" Labels "[\"demo\", \"terraform\"]" FieldName checkbox_example ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Checkbox Example" DefaultChecked true
