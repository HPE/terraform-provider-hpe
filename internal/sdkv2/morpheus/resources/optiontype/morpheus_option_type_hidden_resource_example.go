// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_hidden/resource.tf morpheus_option_type_hidden_resource.tf.tmpl Name tf_example_hidden_option_type Description "Terraform hidden option type example" Labels ["demo","terraform"] FieldName hidden_example ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true DefaultValue example
