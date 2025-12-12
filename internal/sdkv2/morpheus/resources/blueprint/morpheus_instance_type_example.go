// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_instance_type/resource.tf morpheus_instance_type_resource.tf.tmpl Name "tf_example_instance" Code "tf_example_instance" Description "Terraform Example Instance Type" Labels "[\"demo\", \"instance\", \"terraform\"]" Category "web" Visibility "private" ImagePath "tfexample.png" ImageName "tfexample.png" Featured "false" EnableDeployments "true" EnableScaling "true" EnableSettings "true" EnvironmentPrefix "TFEXAMPLE_DEMO" OptionTypeIds "[1910, 1912]" EvarFirstName "first" EvarFirstValue "first" EvarFirstExport "true" EvarSecondName "second" EvarSecondMaskedValue "second" EvarSecondExport "false"
