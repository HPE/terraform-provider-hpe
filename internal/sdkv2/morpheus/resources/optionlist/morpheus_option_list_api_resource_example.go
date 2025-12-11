// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_list_api/resource.tf morpheus_option_list_api_resource.tf.tmpl Name tf_example_api_option_list Description "Terraform Morpheus API option list example" Visibility private OptionList instances TranslationScript "var i=0;\nresults = [];\nfor(i; i<data.length; i++) {\n  results.push({name: data[i].name, value: data[i].name});\n}"
