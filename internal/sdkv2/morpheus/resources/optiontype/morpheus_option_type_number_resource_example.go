// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_number/resource.tf morpheus_option_type_number_resource.tf.tmpl Name tf_example_number_option_type Description "Terraform number option type example" Labels "[\"demo\", \"terraform\"]" FieldName tfNumberExample ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Number Example" Placeholder 12 DefaultValue 1 HelpBlock "Provide a number" Required true
