// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_list_rest/resource.tf morpheus_option_list_rest_resource.tf.tmpl Name tf_example_rest_option_list Description "Terraform REST option list example" Visibility private SourceUrl https://api.github.com/repos/hashicorp/consul/releases RealTime true IgnoreSslErrors true SourceMethod GET InitialDataset "  [{\"name\": \"Level 1\",\"value\":\"level1\"},\n  {\"name\": \"Level 2\",\"value\":\"level2\"},\n  {\"name\": \"Level 3\",\"value\":\"level3\"}\n  ]" TranslationScript "      for(var x=0;x < 5; x++) {\n          results.push({name: data[x].name,value:data[x].name});\n        }" SourceHeaderName1 Accept SourceHeaderValue1 application/json SourceHeaderName2 Authorization SourceHeaderValue2 "Basic YWRtaW46YWRtaW4="
