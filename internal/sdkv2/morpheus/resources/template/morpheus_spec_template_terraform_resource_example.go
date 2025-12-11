// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_terraform/resource_url.tf morpheus_spec_template_terraform_resource_url.tf.tmpl Name "tf-terraform-spec-example-url" SourceType "url" SpecPath "http://example.com/spec.tf"
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_terraform/resource_local.tf morpheus_spec_template_terraform_resource_local.tf.tmpl Name "tf-terraform-spec-example-local" SourceType "local" SpecContent "resource \"aws_instance\" \"instance_1\" {\n  ami           = \"ami-0b91a410940e82c54\"\n  instance_type = \"t2.micro\"\n}"
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_terraform/resource_git.tf morpheus_spec_template_terraform_resource_git.tf.tmpl Name "tf-terraform-spec-example-git" SourceType "repository" RepositoryId "2" VersionRef "main" SpecPath "Instance Types/Terraform/CloudResource/aws/vpc.tf"
