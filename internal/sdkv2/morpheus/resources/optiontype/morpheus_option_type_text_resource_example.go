// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_text/resource.tf morpheus_option_type_text_resource.tf.tmpl Name "tf_example_text_option_type" Description "Terraform text option type example" Labels "[\"demo\", \"terraform\"]" FieldName "test1" ExportMeta "true" DependentField "dependent_example" VisibilityField "visibility_example" RequireField "require_example" ShowOnEdit "true" Editable "true" DisplayValueOnDetails "true" FieldLabel "numbers" Placeholder "fewf" DefaultValue "testing" HelpBlock "fiwefw" Required "true" VerifyPattern "a\\\\D{4}"
