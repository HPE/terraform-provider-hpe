// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_arm/resource_git.tf morpheus_spec_template_arm_resource_git.tf.tmpl Name tf-arm-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./test.json
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_arm/resource_local.tf morpheus_spec_template_arm_resource_local.tf.tmpl Name tf-arm-spec-example-local SourceType local
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_arm/resource_url.tf morpheus_spec_template_arm_resource_url.tf.tmpl Name tf-arm-spec-example-url SourceType url SpecPath http://example.com/spec.json
