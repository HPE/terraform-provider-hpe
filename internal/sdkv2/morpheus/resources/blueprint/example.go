// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_instance_type/resource.tf hpe_morpheus_instance_type_resource.tf.tmpl Name tf_example_instance Code tf_example_instance Description "Terraform Example Instance Type" Labels "[\"demo\", \"instance\", \"terraform\"]" Category web Visibility private ImagePath tfexample.png ImageName tfexample.png Featured false EnableDeployments true EnableScaling true EnableSettings true EnvironmentPrefix TFEXAMPLE_DEMO OptionTypeIds "[1910, 1912]" Evar1Name first Evar1Value first Evar1Export true Evar2Name second Evar2MaskedValue second Evar2Export false
