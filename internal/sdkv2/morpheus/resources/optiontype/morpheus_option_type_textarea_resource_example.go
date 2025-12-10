// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_textarea/resource.tf morpheus_option_type_textarea_resource.tf.tmpl Name tf_example_textarea_option_type Description "Terraform text area option type example" Labels ["demo","terraform"] FieldName textareaExample ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Text Area Example" Rows 5 Placeholder "example text" DefaultValue example HelpBlock "Terraform text area option type example" Required true VerifyPattern "a\\\\D{4}"
