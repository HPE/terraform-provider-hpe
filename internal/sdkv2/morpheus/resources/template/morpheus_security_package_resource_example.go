// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_security_package/resource.tf morpheus_security_package_resource.tf.tmpl Name "tf_example_security_package" Description "Terraform security package example" Labels "[\"demo\", \"terraform\"]" Enabled true Url "https://github.com/ComplianceAsCode/content/releases/download/v0.1.59/scap-security-guide-0.1.59.zip"
