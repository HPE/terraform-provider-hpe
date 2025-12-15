// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_list_manual/resource.tf morpheus_option_list_manual_resource.tf.tmpl Name tf_example_manual_option_list Description "Terraform manual option list example" Dataset "[{\"name\": \"Level 1\",\"value\":\"level1\"},\n {\"name\": \"Level 2\",\"value\":\"level2\"},\n {\"name\": \"Level 3\",\"value\":\"level3\"}\n]" RealTime true
