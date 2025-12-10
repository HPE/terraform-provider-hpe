// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_cluster_package/resource.tf cluster_package_resource.tf.tmpl Name tf_example_cluster_package Code tf-example-cluster-package Description "Terraform example cluster package" PackageVersion 1.2.3 Type apps PackageType example Enabled true RepeatInstall true SpecTemplateIds [1,2]
